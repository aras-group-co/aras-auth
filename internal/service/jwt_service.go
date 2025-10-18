package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
	"github.com/aras-services/aras-auth/pkg/jwt"
)

type JWTService struct {
	jwtService *jwt.JWTService
	tokenRepo  domain.RefreshTokenRepository
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration, tokenRepo domain.RefreshTokenRepository) domain.TokenService {
	return &JWTService{
		jwtService: jwt.NewJWTService(secretKey, accessExpiry, refreshExpiry),
		tokenRepo:  tokenRepo,
	}
}

func (s *JWTService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return s.jwtService.GenerateAccessToken(userID, email)
}

func (s *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	// Generate JWT refresh token
	tokenString, err := s.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return "", err
	}

	// Create refresh token record in database
	tokenID := uuid.New()
	refreshToken := &domain.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		TokenHash: s.hashToken(tokenString),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days default
		CreatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(refreshToken); err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) ValidateAccessToken(token string) (*domain.TokenClaims, error) {
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	return &domain.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		ExpiresAt: claims.ExpiresAt.Unix(),
		IssuedAt:  claims.IssuedAt.Unix(),
		Issuer:    claims.Issuer,
	}, nil
}

func (s *JWTService) ValidateRefreshToken(token string) (*domain.RefreshTokenClaims, error) {
	claims, err := s.jwtService.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	// Verify token exists in database and is not expired
	tokenHash := s.hashToken(token)
	_, err = s.tokenRepo.GetByTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or expired")
	}

	return &domain.RefreshTokenClaims{
		UserID:    claims.UserID,
		TokenID:   claims.TokenID,
		ExpiresAt: claims.ExpiresAt.Unix(),
		IssuedAt:  claims.IssuedAt.Unix(),
		Issuer:    claims.Issuer,
	}, nil
}

func (s *JWTService) RevokeRefreshToken(token string) error {
	// Extract token ID from JWT
	claims, err := s.jwtService.ValidateRefreshToken(token)
	if err != nil {
		return err
	}

	// Delete from database
	return s.tokenRepo.Delete(claims.TokenID)
}

func (s *JWTService) IntrospectToken(token string) (*domain.TokenIntrospection, error) {
	claims, err := s.ValidateAccessToken(token)
	if err != nil {
		return &domain.TokenIntrospection{
			Active: false,
		}, nil
	}

	return &domain.TokenIntrospection{
		Active:    true,
		UserID:    claims.UserID,
		Email:     claims.Email,
		ExpiresAt: claims.ExpiresAt,
		Scope:     "read write", // Default scope for now
	}, nil
}

func (s *JWTService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// CleanupExpiredTokens removes expired refresh tokens from the database
func (s *JWTService) CleanupExpiredTokens() (int, error) {
	return s.tokenRepo.CleanupExpiredTokens()
}

// RevokeAllUserTokens revokes all refresh tokens for a specific user
func (s *JWTService) RevokeAllUserTokens(userID uuid.UUID) error {
	return s.tokenRepo.DeleteByUserID(userID)
}
