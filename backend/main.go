package main

import (
	"context"
	"log/slog"
	"myapp/internal/config"
	"myapp/internal/infrastructure"
	"myapp/internal/logger"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/internal/service/authorization"
	"myapp/internal/transport/handler"
	"myapp/internal/transport/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis_rate/v10"
)

func main() {
	cfg := config.Config{}
	cfg.Load()

	logger.Setup(cfg.LogFormat, cfg.LogLevel)

	slog.Info("starting application")

	//initialize clients
	sqlxClient := infrastructure.NewSQLxClient(cfg.DBConnectionString)
	redisClient := infrastructure.NewRedisClient(cfg.RedisClientHost)

	rateLimiter := redis_rate.NewLimiter(redisClient.GetClient())

	//initialize repositories
	userRepository := repository.NewUserRepository(sqlxClient.GetDB())
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService(cfg.SecretKey)
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(userRepository)

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService, *rateLimiter)
	userHandler := handler.NewUserHandler(userService, *rateLimiter)

	//initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(15 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/auth", authHandler.HandleAuth)
		r.Post("/refresh", authHandler.HandleRefresh)
		r.With(authMiddleware.ProtectionMiddleware).Post("/logout", authHandler.HandleLogout)

		r.Route("/users", func(r chi.Router) {
			r.Post("/create", userHandler.HandleUserCreate)
		})
	})

	httpService := service.NewHTTPService(r)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpService.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Info("server shutdown error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpService.Server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	if err := redisClient.Shutdown(); err != nil {
		slog.Error("redis client shutdown error", "error", err)
	}

	if err := sqlxClient.Shutdown(); err != nil {
		slog.Error("sqlx client shutdown error", "error", err)
	}

	slog.Info("application stopped")
}
