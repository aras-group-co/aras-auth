// Package config implements a centralized configuration management system following
// the 12-Factor App methodology and SOLID principles. It provides type-safe configuration
// loading from multiple sources (YAML files, environment variables, and defaults) with
// proper error handling and graceful degradation.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents the root configuration structure following the Separation of Concerns principle.
// Each field corresponds to a specific functional domain, enabling clear boundaries and
// improved maintainability. The mapstructure tags enable automatic mapping from YAML
// configuration files to Go structs, reducing boilerplate code and providing type safety.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // HTTP server configuration (host, port)
	Database DatabaseConfig `mapstructure:"database"` // Database connection parameters
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT token configuration and expiry settings
	SMTP     SMTPConfig     `mapstructure:"smtp"`     // Email service configuration
	Admin    AdminConfig    `mapstructure:"admin"`    // Default admin user credentials
}

// ServerConfig encapsulates HTTP server configuration following the Single Responsibility Principle.
// Each field uses mapstructure tags for declarative data binding, enabling automatic
// unmarshaling from configuration sources without manual parsing.
type ServerConfig struct {
	Host string `mapstructure:"host"` // Server bind address (default: "0.0.0.0")
	Port int    `mapstructure:"port"` // Server port number (default: 7600)
}

// DatabaseConfig contains PostgreSQL connection parameters using strong typing for
// better compile-time safety and self-documenting code. The SSLMode field allows
// flexible SSL configuration for different deployment environments.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // Database host address
	Port     int    `mapstructure:"port"`     // Database port (default: 5432)
	User     string `mapstructure:"user"`     // Database username
	Password string `mapstructure:"password"` // Database password
	Name     string `mapstructure:"name"`     // Database name
	SSLMode  string `mapstructure:"ssl_mode"` // SSL mode (default: "disable")
}

// JWTConfig manages JWT token settings with strong typing using time.Duration instead
// of strings or integers. This provides compile-time type safety and eliminates
// runtime parsing errors for time-based configurations.
type JWTConfig struct {
	SecretKey     string        `mapstructure:"secret_key"`     // JWT signing secret
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`  // Access token lifetime (default: "15m")
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"` // Refresh token lifetime (default: "7d")
}

// SMTPConfig defines email service configuration for notification and password reset
// functionality. Uses mapstructure tags for consistent configuration mapping patterns.
type SMTPConfig struct {
	Host     string `mapstructure:"host"`     // SMTP server hostname
	Port     int    `mapstructure:"port"`     // SMTP server port
	Username string `mapstructure:"username"` // SMTP authentication username
	Password string `mapstructure:"password"` // SMTP authentication password
	From     string `mapstructure:"from"`     // Default sender email address
}

// AdminConfig stores default administrator credentials for initial system setup.
// This follows the convention over configuration principle by providing sensible defaults.
type AdminConfig struct {
	Email    string `mapstructure:"email"`    // Default admin email address
	Password string `mapstructure:"password"` // Default admin password
}

// Load implements the Configuration Management Pattern with support for multiple configuration sources.
// It follows the 12-Factor App methodology by enabling environment variable overrides and
// providing sensible defaults. The function demonstrates graceful degradation by allowing
// missing configuration files while still failing on actual read errors.
//
// Configuration precedence (highest to lowest):
// 1. Environment variables (via AutomaticEnv)
// 2. Configuration file (config.yaml)
// 3. Default values (fail-safe defaults)
//
// This approach provides flexibility across different deployment environments while
// maintaining type safety and proper error handling.
func Load() (*Config, error) {
	// Configure Viper to look for YAML configuration files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // Look in config directory first
	viper.AddConfigPath(".")        // Fallback to current directory

	// Apply Default Values Pattern: Set fail-safe defaults for optional configuration
	// This follows the "convention over configuration" principle, reducing the burden
	// of configuration while ensuring the application can start with minimal setup.
	viper.SetDefault("server.host", "0.0.0.0")       // Bind to all interfaces by default
	viper.SetDefault("server.port", 7600)            // Standard HTTP port
	viper.SetDefault("database.host", "localhost")   // Local development default
	viper.SetDefault("database.port", 5432)          // PostgreSQL standard port
	viper.SetDefault("database.ssl_mode", "disable") // Disable SSL for local development
	viper.SetDefault("jwt.access_expiry", "15m")     // Short-lived access tokens for security
	viper.SetDefault("jwt.refresh_expiry", "7d")     // Longer-lived refresh tokens

	// Enable 12-Factor App compliance: Environment variables override file configuration
	// This allows deployment-specific configuration without modifying code or files
	viper.AutomaticEnv()

	// Implement Graceful Degradation: Allow missing config file, fail only on read errors
	// This enables environment-variable-only configurations for containerized deployments
	if err := viper.ReadInConfig(); err != nil {
		// Only fail if there's an actual read error, not if the file is missing
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal configuration into strongly-typed struct using mapstructure tags
	// This provides compile-time type safety and eliminates runtime parsing errors
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
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
