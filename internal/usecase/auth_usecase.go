package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
	"github.com/aras-services/aras-auth/pkg/password"
)

type AuthUseCase struct {
	providerRegistry domain.ProviderRegistry
	tokenService     domain.TokenService
	userRepo         domain.UserRepository
}

func NewAuthUseCase(providerRegistry domain.ProviderRegistry, tokenService domain.TokenService, userRepo domain.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		providerRegistry: providerRegistry,
		tokenService:     tokenService,
		userRepo:         userRepo,
	}
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	TokenType    string       `json:"token_type"`
	User         *domain.User `json:"user"`
}

type RegisterResponse struct {
	User    *domain.User `json:"user"`
	Message string       `json:"message"`
}

func (uc *AuthUseCase) Register(ctx context.Context, req *domain.CreateUserRequest) (*RegisterResponse, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Validate password
	if !password.IsValidPassword(req.Password) {
		return nil, fmt.Errorf("password does not meet requirements")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Status:        domain.UserStatusPending,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Get default provider and create user
	provider := uc.providerRegistry.GetDefaultProvider()
	if provider == nil {
		return nil, fmt.Errorf("no identity provider available")
	}

	if err := provider.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &RegisterResponse{
		User:    user,
		Message: "User registered successfully. Please verify your email.",
	}, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*LoginResponse, error) {
	// Get default provider
	provider := uc.providerRegistry.GetDefaultProvider()
	if provider == nil {
		return nil, fmt.Errorf("no identity provider available")
	}

	// Authenticate user
	user, err := provider.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token
	claims, err := uc.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := uc.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is still active
	if user.Status != domain.UserStatusActive {
		return nil, fmt.Errorf("user account is not active")
	}

	// Generate new access token
	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate new refresh token (token rotation)
	newRefreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Revoke old refresh token
	if err := uc.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: failed to revoke old refresh token: %v\n", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900, // 15 minutes
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	// Revoke refresh token
	return uc.tokenService.RevokeRefreshToken(refreshToken)
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error {
	// Verify current password
	provider := uc.providerRegistry.GetDefaultProvider()
	if provider == nil {
		return fmt.Errorf("no identity provider available")
	}

	valid, err := provider.VerifyPassword(ctx, userID, req.CurrentPassword)
	if err != nil || !valid {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if !password.IsValidPassword(req.NewPassword) {
		return fmt.Errorf("new password does not meet requirements")
	}

	// Change password
	return provider.ChangePassword(ctx, userID, req.NewPassword)
}

func (uc *AuthUseCase) ForgotPassword(ctx context.Context, req *domain.ResetPasswordRequest) error {
	// Check if user exists
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil
	}

	// TODO: Send password reset email
	// For now, just log the request
	fmt.Printf("Password reset requested for user: %s\n", user.Email)

	return nil
}

func (uc *AuthUseCase) ResetPassword(ctx context.Context, req *domain.ConfirmResetPasswordRequest) error {
	// TODO: Validate reset token
	// For now, this is a placeholder implementation

	// Find user by token (in a real implementation, you'd store tokens in DB)
	// For now, return an error
	return fmt.Errorf("password reset not implemented yet")
}

func (uc *AuthUseCase) VerifyEmail(ctx context.Context, userID uuid.UUID) error {
	// Update user email verification status
	return uc.userRepo.UpdateEmailVerified(userID, true)
}

func (uc *AuthUseCase) IntrospectToken(ctx context.Context, token string) (*domain.TokenIntrospection, error) {
	return uc.tokenService.IntrospectToken(token)
}
