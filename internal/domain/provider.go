package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Provider represents an identity provider configuration
type Provider struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	Type      string          `json:"type" db:"type"`
	Config    json.RawMessage `json:"config" db:"config"`
	Enabled   bool            `json:"enabled" db:"enabled"`
	IsSystem  bool            `json:"is_system" db:"is_system"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// IdentityProvider defines the interface for identity providers
// This allows for pluggable authentication backends
type IdentityProvider interface {
	// Authenticate verifies user credentials and returns user info
	Authenticate(ctx context.Context, username, password string) (*User, error)

	// GetUser retrieves a user by ID
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *User) error

	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error

	// VerifyPassword verifies a password against user's hash
	VerifyPassword(ctx context.Context, userID uuid.UUID, password string) (bool, error)

	// GetProviderName returns the name of this provider
	GetProviderName() string

	// IsEnabled returns whether this provider is currently enabled
	IsEnabled() bool
}

// ProviderRegistry manages multiple identity providers
type ProviderRegistry interface {
	// RegisterProvider registers a new identity provider
	RegisterProvider(provider IdentityProvider) error

	// GetProvider retrieves a provider by name
	GetProvider(name string) (IdentityProvider, error)

	// GetDefaultProvider returns the default provider
	GetDefaultProvider() IdentityProvider

	// ListProviders returns all registered providers
	ListProviders() []IdentityProvider

	// GetEnabledProviders returns only enabled providers
	GetEnabledProviders() []IdentityProvider
}

// TokenService handles JWT token operations
type TokenService interface {
	// GenerateAccessToken creates a new access token for a user
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)

	// GenerateRefreshToken creates a new refresh token for a user
	GenerateRefreshToken(userID uuid.UUID) (string, error)

	// ValidateAccessToken validates an access token and returns claims
	ValidateAccessToken(token string) (*TokenClaims, error)

	// ValidateRefreshToken validates a refresh token
	ValidateRefreshToken(token string) (*RefreshTokenClaims, error)

	// RevokeRefreshToken invalidates a refresh token
	RevokeRefreshToken(token string) error

	// IntrospectToken provides token information for other services
	IntrospectToken(token string) (*TokenIntrospection, error)
}

// TokenClaims represents the claims in an access token
type TokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	ExpiresAt int64     `json:"exp"`
	IssuedAt  int64     `json:"iat"`
	Issuer    string    `json:"iss"`
}

// RefreshTokenClaims represents the claims in a refresh token
type RefreshTokenClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenID   uuid.UUID `json:"token_id"`
	ExpiresAt int64     `json:"exp"`
	IssuedAt  int64     `json:"iat"`
	Issuer    string    `json:"iss"`
}

// TokenIntrospection provides token information for external services
type TokenIntrospection struct {
	Active    bool      `json:"active"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	ExpiresAt int64     `json:"exp,omitempty"`
	Scope     string    `json:"scope,omitempty"`
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RefreshTokenRepository handles refresh token persistence
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByID(id uuid.UUID) (*RefreshToken, error)
	GetByTokenHash(tokenHash string) (*RefreshToken, error)
	GetByUserID(userID uuid.UUID) ([]*RefreshToken, error)
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
	CleanupExpiredTokens() (int, error)
}
