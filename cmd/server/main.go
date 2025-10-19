// Package main implements the server entry point for the ARAS Authentication Service.
// This application follows Clean Architecture principles with clear separation of concerns
// across multiple layers: Repository (data access) → UseCase (business logic) → Handler (HTTP interface).
// The main function demonstrates Dependency Injection, Factory patterns, and graceful shutdown
// handling, showcasing enterprise-grade Go application architecture.
package main

import (
	"context"
	"fmt"
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

// Version information - set during build time via ldflags
var (
	version   = "1.1.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// printVersion prints version information and exits
func printVersion() {
	fmt.Printf("aras_auth version %s\n", version)
	if buildTime != "unknown" {
		fmt.Printf("Build Time: %s\n", buildTime)
	}
	if gitCommit != "unknown" {
		fmt.Printf("Git Commit: %s\n", gitCommit)
	}
	os.Exit(0)
}

// main implements the application bootstrap following Clean Architecture principles.
// It demonstrates Dependency Injection, Factory patterns, and proper resource management
// while maintaining clear separation of concerns across architectural layers.
func main() {


	// Check for version flag before any initialization
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if arg == "--version" || arg == "-v" {
				printVersion()
			}
		}
	}

	// PHASE 1: Configuration and Infrastructure Setup
	// Load configuration using the centralized config management pattern
	// This follows the 12-Factor App methodology for configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize structured logger using Factory pattern
	// Zap provides high-performance structured logging with minimal allocation overhead
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync() // Resource Management: Ensure log buffer is flushed on exit

	// PHASE 2: Database Connection and Health Check
	// Connect to PostgreSQL using connection pooling for optimal performance
	// The DSN is constructed using the config's encapsulated helper method
	db, err := pgxpool.New(context.Background(), cfg.GetDSN())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close() // Resource Management: Ensure database connection is closed

	// Test database connectivity using fail-fast pattern
	// This ensures the application fails early if database is unreachable
	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Connected to database successfully")

	// PHASE 3: Repository Layer Initialization (Data Access Layer)
	// Repository Pattern: Abstract data access through interfaces
	// Each repository encapsulates database operations for a specific domain entity
	// This follows the Single Responsibility Principle and enables easy testing
	userRepo := postgres.NewUserRepository(db)             // Factory Pattern: Constructor injection
	groupRepo := postgres.NewGroupRepository(db)           // Concrete PostgreSQL implementation
	roleRepo := postgres.NewRoleRepository(db)             // Implements domain interfaces
	permissionRepo := postgres.NewPermissionRepository(db) // Dependency Inversion Principle
	tokenRepo := postgres.NewTokenRepository(db)           // Depends on abstractions, not concrete types

	// PHASE 4: Service Layer Initialization (Cross-cutting Concerns)
	// JWT Service handles token generation and validation
	// Uses constructor injection with configuration and repository dependencies
	jwtService := service.NewJWTService(
		cfg.JWT.SecretKey,     // Configuration injection
		cfg.JWT.AccessExpiry,  // Type-safe duration from config
		cfg.JWT.RefreshExpiry, // Follows Dependency Injection pattern
		tokenRepo,             // Repository dependency injection
	)

	// PHASE 5: Provider Registry Pattern (Plugin Architecture)
	// Registry Pattern: Manages pluggable authentication providers
	// Enables Open/Closed Principle - open for extension, closed for modification
	// Supports multiple identity providers (local, OAuth, LDAP, SAML, etc.)
	providerRegistry := provider.NewProviderRegistry()

	// Register local authentication provider
	// Adapter Pattern: Adapts local user repository to provider interface
	localProvider := local.NewLocalProvider(userRepo)
	if err := providerRegistry.RegisterProvider(localProvider); err != nil {
		logger.Fatal("Failed to register local provider", zap.Error(err))
	}

	// PHASE 6: Use Case Layer Initialization (Business Logic Layer)
	// Use Case Pattern: Encapsulates business logic and orchestrates domain operations
	// Each use case handles a specific business capability and coordinates between
	// repositories, services, and external dependencies
	authUseCase := usecase.NewAuthUseCase(providerRegistry, jwtService, userRepo) // Authentication business logic
	userUseCase := usecase.NewUserUseCase(userRepo)                               // User management business logic
	groupUseCase := usecase.NewGroupUseCase(groupRepo)                            // Group management business logic
	authzUseCase := usecase.NewAuthzUseCase(roleRepo, permissionRepo)             // Authorization business logic

	// PHASE 7: Handler Layer Initialization (Interface Adapters)
	// Adapter Pattern: HTTP handlers adapt external HTTP requests to use cases
	// Each handler is responsible for HTTP-specific concerns (parsing, validation, response formatting)
	// while delegating business logic to use cases
	authHandler := httphandler.NewAuthHandler(authUseCase)    // Authentication HTTP interface
	userHandler := httphandler.NewUserHandler(userUseCase)    // User management HTTP interface
	groupHandler := httphandler.NewGroupHandler(groupUseCase) // Group management HTTP interface
	authzHandler := httphandler.NewAuthzHandler(authzUseCase) // Authorization HTTP interface

	// PHASE 8: Middleware Initialization (Cross-cutting Concerns)
	// Middleware Pattern: Chain of Responsibility for cross-cutting concerns
	// Each middleware wraps handlers with additional behavior (auth, logging, CORS, etc.)
	authMiddleware := authmiddleware.NewAuthMiddleware(jwtService)     // JWT token validation
	rbacMiddleware := authmiddleware.NewRBACMiddleware(permissionRepo) // Role-based access control
	corsMiddleware := authmiddleware.NewCORSMiddleware()               // Cross-origin resource sharing

	// PHASE 9: Router Configuration and Middleware Chain Setup
	// Router Pattern: Hierarchical route organization with middleware scoping
	// Chi router provides lightweight, idiomatic HTTP routing with middleware support
	r := chi.NewRouter()

	// Middleware Chain: Chain of Responsibility pattern
	// Middleware is applied in order - each wraps the next handler in the chain
	// This enables composable cross-cutting concerns with clear separation
	r.Use(middleware.Logger)                    // Request logging middleware
	r.Use(middleware.Recoverer)                 // Panic recovery middleware
	r.Use(corsMiddleware)                       // CORS handling middleware
	r.Use(middleware.RequestID)                 // Request ID generation for tracing
	r.Use(middleware.RealIP)                    // Real IP extraction (behind proxies)
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout protection

	// Health Check Endpoint Pattern
	// Dedicated health endpoint for monitoring and orchestration systems (Kubernetes, Docker)
	// Provides immediate feedback on service availability
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// PHASE 10: API Route Configuration with Layered Security
	// Route Grouping Pattern: Hierarchical route organization with middleware scoping
	// Routes are organized by functionality with appropriate security boundaries
	r.Route("/api/v1", func(r chi.Router) {
		// Public Routes: No authentication required
		// Authentication endpoints (login, register, password reset) are publicly accessible
		authHandler.RegisterRoutes(r)

		// Protected Routes: Require authentication
		// Route Group Pattern: Scoped middleware application
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth) // Authentication middleware applied to all routes in this group

			// User Management Routes: Authenticated users can manage their own data
			userHandler.RegisterRoutes(r)

			// Group Management Routes: Require specific permissions
			// Nested Route Groups: Fine-grained permission control
			r.Group(func(r chi.Router) {
				r.Use(rbacMiddleware.RequirePermission("groups", "read")) // RBAC middleware for group operations
				groupHandler.RegisterRoutes(r)
			})

			// Authorization Routes: Require admin-level permissions
			// Role and permission management requires elevated privileges
			r.Group(func(r chi.Router) {
				r.Use(rbacMiddleware.RequirePermission("roles", "read")) // RBAC middleware for role operations
				authzHandler.RegisterRoutes(r)
			})
		})
	})

	// PHASE 11: Server Initialization and Startup
	// HTTP Server Configuration: Using Go's standard http.Server with custom handler
	// Server address is constructed using the config's encapsulated helper method
	server := &http.Server{
		Addr:    cfg.GetServerAddr(), // Configuration-driven server binding
		Handler: r,                   // Chi router as the main handler
	}

	// Concurrent Server Startup Pattern
	// Start server in a goroutine to allow main thread to handle shutdown signals
	// This enables graceful shutdown without blocking the startup process
	go func() {
		logger.Info("Starting server", zap.String("addr", cfg.GetServerAddr()))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// PHASE 12: Graceful Shutdown Implementation
	// Signal Handling Pattern: Listen for OS signals to initiate graceful shutdown
	// This ensures clean resource cleanup and prevents dropped requests
	quit := make(chan os.Signal, 1)                      // Buffered channel for signal handling
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Register for interrupt and terminate signals
	<-quit                                               // Block until shutdown signal received

	logger.Info("Shutting down server...")

	// Graceful Shutdown with Timeout Pattern
	// Create a deadline context to prevent indefinite shutdown waiting
	// This ensures the server doesn't hang indefinitely on shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Ensure context is cancelled to prevent resource leaks

	// Attempt graceful shutdown with timeout
	// This allows in-flight requests to complete while preventing indefinite waiting
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
