package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/aras-services/aras-auth/config"
	httphandler "github.com/aras-services/aras-auth/internal/delivery/http"
	authmiddleware "github.com/aras-services/aras-auth/internal/middleware"
	"github.com/aras-services/aras-auth/internal/provider"
	"github.com/aras-services/aras-auth/internal/provider/local"
	"github.com/aras-services/aras-auth/internal/repository/postgres"
	"github.com/aras-services/aras-auth/internal/service"
	"github.com/aras-services/aras-auth/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Connect to database
	db, err := pgxpool.New(context.Background(), cfg.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Connected to database successfully")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	groupRepo := postgres.NewGroupRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permissionRepo := postgres.NewPermissionRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)

	// Initialize JWT service
	jwtService := service.NewJWTService(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessExpiry,
		cfg.JWT.RefreshExpiry,
		tokenRepo,
	)

	// Initialize provider registry
	providerRegistry := provider.NewProviderRegistry()

	// Register local provider
	localProvider := local.NewLocalProvider(userRepo)
	if err := providerRegistry.RegisterProvider(localProvider); err != nil {
		logger.Fatal("Failed to register local provider", zap.Error(err))
	}

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(providerRegistry, jwtService, userRepo)
	userUseCase := usecase.NewUserUseCase(userRepo)
	groupUseCase := usecase.NewGroupUseCase(groupRepo)
	authzUseCase := usecase.NewAuthzUseCase(roleRepo, permissionRepo)

	// Initialize handlers
	authHandler := httphandler.NewAuthHandler(authUseCase)
	userHandler := httphandler.NewUserHandler(userUseCase)
	groupHandler := httphandler.NewGroupHandler(groupUseCase)
	authzHandler := httphandler.NewAuthzHandler(authzUseCase)

	// Initialize middleware
	authMiddleware := authmiddleware.NewAuthMiddleware(jwtService)
	rbacMiddleware := authmiddleware.NewRBACMiddleware(permissionRepo)
	corsMiddleware := authmiddleware.NewCORSMiddleware()

	// Setup router
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth required)
		authHandler.RegisterRoutes(r)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			// User management routes
			userHandler.RegisterRoutes(r)

			// Group management routes (require admin permissions)
			r.Group(func(r chi.Router) {
				r.Use(rbacMiddleware.RequirePermission("groups", "read"))
				groupHandler.RegisterRoutes(r)
			})

			// Authorization routes (require admin permissions)
			r.Group(func(r chi.Router) {
				r.Use(rbacMiddleware.RequirePermission("roles", "read"))
				authzHandler.RegisterRoutes(r)
			})
		})
	})

	// Start server
	server := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("addr", cfg.GetServerAddr()))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
