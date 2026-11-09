package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/handler"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/jwt"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/middleware"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/proxy"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/repo"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/adapter/tripclient"
	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/usecase"
)

type realClock struct{}

func (realClock) Now() time.Time { return time.Now().UTC() }

type uuidGen struct{}

func (uuidGen) New() string { return uuid.NewString() }

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "atlasflow-default-secret-change-me-32-bytes-long!")
	tripServiceURL := getEnv("TRIP_SERVICE_URL", "http://localhost:8081")
	geoServiceURL := getEnv("GEO_SERVICE_URL", "http://localhost:8082")
	dbURL := os.Getenv("DATABASE_URL")

	clock := realClock{}
	idGen := uuidGen{}

	var userRepo usecase.UserRepository
	var tokenRepo usecase.RefreshTokenRepository

	if dbURL != "" {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			logger.Error("failed to open database, falling back to in-memory", "error", err)
			userRepo = repo.NewInMemoryUserRepository()
			tokenRepo = repo.NewInMemoryTokenRepository()
		} else {
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(10)
			db.SetConnMaxLifetime(5 * time.Minute)
			userRepo = repo.NewPostgresUserRepository(db)
			tokenRepo = repo.NewPostgresTokenRepository(db)
			logger.Info("connected to PostgreSQL database")
		}
	} else {
		logger.Info("DATABASE_URL not set; using in-memory thread-safe repositories")
		userRepo = repo.NewInMemoryUserRepository()
		tokenRepo = repo.NewInMemoryTokenRepository()
	}

	tokenManager := jwt.NewTokenManager(jwtSecret, 15*time.Minute, 7*24*time.Hour, clock)
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenRepo, tokenManager, clock, idGen)

	tripClient := tripclient.NewHTTPTripClient(tripServiceURL)
	shareUsecase := usecase.NewPublicShareUsecase(tripClient)

	authHandler := handler.NewAuthHandler(authUsecase)
	shareHandler := handler.NewShareHandler(shareUsecase)
	authMiddleware := middleware.NewAuthMiddleware(tokenManager)

	rps, _ := strconv.ParseFloat(getEnv("RATE_LIMIT_RPS", "20"), 64)
	burst, _ := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "50"))
	rateLimiter := middleware.NewIPRateLimiter(rps, burst)

	tripProxy, err := proxy.NewReverseProxy(tripServiceURL, logger)
	if err != nil {
		logger.Error("failed to initialize trip reverse proxy", "error", err)
	}
	geoProxy, err := proxy.NewReverseProxy(geoServiceURL, logger)
	if err != nil {
		logger.Error("failed to initialize geo reverse proxy", "error", err)
	}

	mux := http.NewServeMux()

	// 1. Healthcheck
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"gateway-auth"}`))
	})

	// 2. ConnectRPC AuthService procedure endpoints
	mux.Handle("/atlasflow.auth.v1.AuthService/", authMiddleware.OptionalAuth(http.HandlerFunc(authHandler.ServeConnectRPC)))

	// 3. Public Auth REST endpoints
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.RegisterHTTP)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.LoginHTTP)
	mux.HandleFunc("POST /api/v1/auth/refresh", authHandler.RefreshHTTP)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.LogoutHTTP)

	// 4. Protected Auth REST endpoint
	mux.Handle("GET /api/v1/auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.MeHTTP)))

	// 5. Public Share endpoints (no JWT, rate limited, strictly sanitized)
	shareRoute := rateLimiter.Limit(http.HandlerFunc(shareHandler.GetShareTrip))
	mux.Handle("GET /api/v1/share/{slug}", shareRoute)
	mux.Handle("GET /share/{slug}", shareRoute)

	// 6. Internal Service Proxies (protected, forwards X-User-Id)
	if tripProxy != nil {
		mux.Handle("/api/v1/trips/", authMiddleware.RequireAuth(tripProxy))
	}
	if geoProxy != nil {
		mux.Handle("/api/v1/geo/", authMiddleware.RequireAuth(geoProxy))
	}

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version, X-User-Agent")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		mux.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("AtlasFlow Gateway & Auth service running", "port", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	sig := <-shutdown
	logger.Info("shutting down gateway gracefully", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "error", err)
		_ = server.Close()
	}

	logger.Info("gateway server cleanly stopped")
}