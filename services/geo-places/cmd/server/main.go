package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/geocode"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/handler"
	"github.com/makskozlov/atlasflow/services/geo-places/adapter/repo"
	"github.com/makskozlov/atlasflow/services/geo-places/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	dbURL := os.Getenv("DATABASE_URL")
	nominatimURL := os.Getenv("NOMINATIM_URL")

	var placeRepo usecase.PlaceRepository
	if dbURL != "" {
		logger.Info("connecting to postgis database")
		db, err := sql.Open("pgx", dbURL)
		if err != nil {
			logger.Error("failed to open database connection", "error", err)
			os.Exit(1)
		}
		defer db.Close()

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(15 * time.Minute)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := db.PingContext(ctx); err != nil {
			logger.Warn("database ping failed, falling back to memory repository", "error", err)
			cancel()
			placeRepo = repo.NewMemoryPlaceRepository()
		} else {
			cancel()
			logger.Info("connected to postgis database successfully")
			placeRepo = repo.NewPostgresPlaceRepository(db)
		}
	} else {
		logger.Info("DATABASE_URL not set, using in-memory repository")
		placeRepo = repo.NewMemoryPlaceRepository()
	}

	var geocoder usecase.Geocoder
	if nominatimURL != "" {
		logger.Info("using Nominatim geocoder", "url", nominatimURL)
		geocoder = geocode.NewNominatimGeocoder(nominatimURL, nil)
	} else {
		logger.Info("NOMINATIM_URL not set, using mock geocoder")
		geocoder = geocode.NewMockGeocoder()
	}

	geoUsecase := usecase.NewGeoUsecase(placeRepo, geocoder, nil)
	geoHandler := handler.NewGeoHandler(geoUsecase)
	httpHandler := handler.NewHTTPHandler(geoHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      httpHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("geo-places service starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	logger.Info("shutting down geo-places server gracefully...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}
	logger.Info("geo-places server exited cleanly")
}