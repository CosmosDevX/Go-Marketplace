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
	ctx := context.Background()
	cfg := config.Config{}
	cfg.Load()

	logger.Setup(cfg.LogFormat, cfg.LogLevel)
	slog.Info("starting application")

	//initialize clients
	sqlxClient := infrastructure.NewSQLxClient(cfg.DBConnectionString)
	redisClient := infrastructure.NewRedisClient(ctx, cfg.RedisClientHost, cfg.RedisClientPassword)
	rateLimiter := redis_rate.NewLimiter(redisClient.GetClient())

	//initialize repositories
	userRepository := repository.NewUserRepository(sqlxClient.GetDB())
	categoryRepository := repository.NewCategoryRepository(sqlxClient.GetDB())
	productRepository := repository.NewProductRepository(sqlxClient.GetDB())
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService(cfg.SecretKey)
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(userRepository)
	categoryService := service.NewCategoryService(categoryRepository)
	productService := service.NewProductService(productRepository, categoryRepository)

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService, *rateLimiter)
	userHandler := handler.NewUserHandler(userService, *rateLimiter)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)

	//initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	//constants
	const maxBodySize = 1024 * 1024

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.MaxBodySize(maxBodySize))
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(15 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth", authHandler.Auth)
		r.Post("/refresh", authHandler.Refresh)
		r.With(authMiddleware.ProtectionMiddleware).Post("/logout", authHandler.Logout)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateHandler)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Post("/", categoryHandler.CreateHandler)
			r.Get("/", categoryHandler.ListHandler)
		})

		r.Route("/products", func(r chi.Router) {
			r.Post("/", productHandler.CreateHandler)
			r.Get("/", productHandler.ListHandler)
		})

		r.Get("/uploads/{file}", func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/api/v1/uploads/", http.FileServer(http.Dir("./uploads"))).ServeHTTP(w, r)
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
