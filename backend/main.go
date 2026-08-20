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
	"myapp/internal/utils"
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

	fileManager := utils.NewFileManager()

	logger.Setup(cfg.LogFormat, cfg.LogLevel)
	slog.Info("starting application")

	//initialize clients
	sqlxClient := infrastructure.NewSQLxClient(cfg.DBConnectionString)
	redisClient := infrastructure.NewRedisClient(ctx, cfg.RedisClientHost, cfg.RedisClientPassword)
	rateLimiter := redis_rate.NewLimiter(redisClient.GetClient())

	unitOfWork := service.NewUnitOfWork(sqlxClient.GetDB())

	//initialize repositories
	userRepository := repository.NewUserRepository(sqlxClient.GetDB())
	userRoleRepository := repository.NewUserRoleRepository(sqlxClient.GetDB())
	categoryRepository := repository.NewCategoryRepository(sqlxClient.GetDB())
	productRepository := repository.NewProductRepository(sqlxClient.GetDB())
	cartRepository := repository.NewCartRepository(sqlxClient.GetDB())
	cartItemRepository := repository.NewCartItemRepository(sqlxClient.GetDB())
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService(cfg.SecretKey)
	authService := authorization.NewAuthService(userRepository, userRoleRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(unitOfWork)
	categoryService := service.NewCategoryService(categoryRepository)
	productService := service.NewProductService(unitOfWork, productRepository, categoryRepository, fileManager)
	cartItemService := service.NewCartItemService(cartItemRepository, cartRepository)

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService, *rateLimiter)
	userHandler := handler.NewUserHandler(userService, *rateLimiter)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	cartItemHandler := handler.NewCartItemHandler(cartItemService)

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
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.Create)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryHandler.List)
			r.With(authMiddleware.ProtectionMiddleware, middleware.RoleMiddleware([]string{config.AdminRole})).Post("/", categoryHandler.Create)
		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.List)
		})

		r.Route("/seller/products", func(r chi.Router) {
			r.Use(authMiddleware.ProtectionMiddleware, middleware.RoleMiddleware([]string{config.SellerRole, config.AdminRole}))
			r.Get("/", productHandler.ListBySeller)
			r.Post("/", productHandler.Create)
			r.Delete("/{product_id}", productHandler.Delete)
		})

		r.Route("/cart", func(r chi.Router) {
			r.With(authMiddleware.ProtectionMiddleware).Get("/", cartItemHandler.List)
			r.With(authMiddleware.ProtectionMiddleware).Post("/items", cartItemHandler.Create)
			r.With(authMiddleware.ProtectionMiddleware).Patch("/items/{cart_item_id}", cartItemHandler.UpdateQuantity)
			r.With(authMiddleware.ProtectionMiddleware).Delete("/items/{cart_item_id}", cartItemHandler.Delete)
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
