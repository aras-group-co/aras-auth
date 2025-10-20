// Package config implements a centralized configuration management system following
// the 12-Factor App methodology and SOLID principles. It provides type-safe configuration
// loading from multiple sources (YAML files, environment variables, and defaults) with
// proper error handling and graceful degradation.
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config represents the root configuration structure following the Separation of Concerns principle.
// Each field corresponds to a specific functional domain, enabling clear boundaries and
// improved maintainability. The env tags enable automatic mapping from environment variables
// to Go structs, reducing boilerplate code and providing type safety.
type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Database DatabaseConfig `envPrefix:"DB_"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	SMTP     SMTPConfig     `envPrefix:"SMTP_"`
	Admin    AdminConfig    `envPrefix:"ADMIN_"`
}

// ServerConfig encapsulates HTTP server configuration following the Single Responsibility Principle.
// Each field uses env tags for declarative data binding, enabling automatic
// unmarshaling from environment variables without manual parsing.
type ServerConfig struct {
	Host string `env:"HOST" envDefault:"0.0.0.0"` // Server bind address (default: "0.0.0.0")
	Port int    `env:"PORT" envDefault:"7600"`    // Server port number (default: 7600)
}

// DatabaseConfig contains PostgreSQL connection parameters using strong typing for
// better compile-time safety and self-documenting code. The SSLMode field allows
// flexible SSL configuration for different deployment environments.
type DatabaseConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`    // Database host address
	Port     int    `env:"PORT" envDefault:"5432"`         // Database port (default: 5432)
	User     string `env:"USER" envDefault:"postgres"`     // Database username
	Password string `env:"PASSWORD" envDefault:"postgres"` // Database password
	Name     string `env:"NAME" envDefault:"aras_auth"`    // Database name
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`  // SSL mode (default: "disable")
}

// JWTConfig manages JWT token settings with strong typing using time.Duration instead
// of strings or integers. This provides compile-time type safety and eliminates
// runtime parsing errors for time-based configurations.
type JWTConfig struct {
	SecretKey     string        `env:"SECRET_KEY" envDefault:"change-me-please-32b-min"` // JWT signing secret
	AccessExpiry  time.Duration `env:"ACCESS_EXPIRY" envDefault:"15m"`                   // Access token lifetime (default: "15m")
	RefreshExpiry time.Duration `env:"REFRESH_EXPIRY" envDefault:"168h"`                 // Refresh token lifetime (default: "168h")
}

// SMTPConfig defines email service configuration for notification and password reset
// functionality. Uses env tags for consistent configuration mapping patterns.
type SMTPConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`                 // SMTP server hostname
	Port     int    `env:"PORT" envDefault:"587"`                       // SMTP server port
	Username string `env:"USERNAME" envDefault:""`                      // SMTP authentication username
	Password string `env:"PASSWORD" envDefault:""`                      // SMTP authentication password
	From     string `env:"FROM" envDefault:"noreply@aras-services.com"` // Default sender email address
}

// AdminConfig stores default administrator credentials for initial system setup.
// This follows the convention over configuration principle by providing sensible defaults.
type AdminConfig struct {
	Email    string `env:"EMAIL" envDefault:"admin@aras-services.com"` // Default admin email address
	Password string `env:"PASSWORD" envDefault:"admin123"`             // Default admin password
}

// Load implements the Configuration Management Pattern with support for environment variables only.
// It follows the 12-Factor App methodology by reading all configuration from environment variables
// with sensible defaults. This approach provides maximum flexibility across different deployment
// environments while maintaining type safety and proper error handling.
//
// Configuration precedence (highest to lowest):
// 1. Environment variables
// 2. Default values (fail-safe defaults)
//
// This approach provides flexibility across different deployment environments while
// maintaining type safety and proper error handling.
func Load() (*Config, error) {
	var config Config

	// Parse environment variables into the config struct
	// This provides compile-time type safety and eliminates runtime parsing errors
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return &config, nil
}

// GetDSN implements the Encapsulation pattern by providing a single method to construct
// the PostgreSQL Data Source Name (DSN). This follows the DRY principle by centralizing
// the connection string formatting logic, making it easier to modify the format if needed
// and ensuring consistency across the application.
//
// Benefits:
// - Single source of truth for database connection string format
// - Easier to modify connection string format in the future
// - Reduces code duplication across different parts of the application
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetServerAddr implements the Encapsulation pattern by providing a centralized method
// to construct the server address string. This encapsulates the string formatting logic
// and provides a clean interface for obtaining the server's bind address.
//
// Benefits:
// - Single source of truth for server address format
// - Easier to modify address format if needed (e.g., adding IPv6 support)
// - Consistent interface for server address construction
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
