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

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/adapter/handler"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/adapter/repo"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/adapter/storage"
	"github.com/makskozlov/atlasflow/services/trip-core/internal/usecase"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type realClock struct{}

func (realClock) Now() time.Time { return time.Now().UTC() }

type uuidGen struct{}

func (uuidGen) New() string { return uuid.NewString() }

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := getEnv("PORT", "8081")
	dbURL := os.Getenv("DATABASE_URL")
	s3Endpoint := getEnv("S3_ENDPOINT", "http://localhost:9000")
	s3PublicEndpoint := getEnv("S3_PUBLIC_ENDPOINT", s3Endpoint)
	s3Bucket := getEnv("S3_BUCKET", "atlasflow-media")

	clock := realClock{}
	ids := uuidGen{}

	var tripRepo usecase.TripRepository
	if dbURL != "" {
		logger.Info("connecting to postgresql database", "url", dbURL)
		db, err := sql.Open("pgx", dbURL)
		if err != nil {
			logger.Error("failed to open database connection, falling back to memory", "error", err)
			tripRepo = repo.NewMemoryTripRepo()
		} else {
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(5)
			db.SetConnMaxLifetime(15 * time.Minute)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := db.PingContext(ctx); err != nil {
				logger.Warn("database ping failed, using memory repository", "error", err)
				cancel()
				tripRepo = repo.NewMemoryTripRepo()
			} else {
				cancel()
				logger.Info("connected to PostgreSQL successfully")
				tripRepo = repo.NewPostgresTripRepo(db)
			}
		}
	} else {
		logger.Info("DATABASE_URL not set; using thread-safe in-memory repository")
		tripRepo = repo.NewMemoryTripRepo()
	}

	storageAdapter := storage.NewS3Storage(s3Endpoint, s3PublicEndpoint, s3Bucket)

	tu := usecase.NewTripUsecase(tripRepo, clock, ids)
	du := usecase.NewDayUsecase(tripRepo, clock, ids)
	iu := usecase.NewItineraryUsecase(tripRepo, clock, ids)
	bu := usecase.NewBudgetUsecase(tripRepo, clock, ids)
	cu := usecase.NewChecklistUsecase(tripRepo, clock, ids)
	au := usecase.NewAttachmentUsecase(tripRepo, storageAdapter, clock, ids)
	su := usecase.NewPublicShareUsecase(tripRepo, clock, ids)

	tripHandler := handler.NewTripHandler(tu, du, iu, bu, cu, au, su)

	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"trip-core"}`))
	})

	// Forward X-User-Id header into context for internal requests
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if uid := r.Header.Get("X-User-Id"); uid != "" {
				ctx = handler.ContextWithUserID(ctx, uid)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	httpHandler := handler.NewHTTPHandler(tripHandler)
	mux.Handle("/atlasflow.trip.v1.TripService/", authMiddleware(httpHandler))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("AtlasFlow Trip-Core service running", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-shutdown
	logger.Info("shutting down trip-core gracefully", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		_ = server.Close()
	}

	logger.Info("trip-core cleanly stopped")
}
