package local

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
	"github.com/aras-services/aras-auth/pkg/password"
)

type LocalProvider struct {
	userRepo domain.UserRepository
}

func NewLocalProvider(userRepo domain.UserRepository) domain.IdentityProvider {
	return &LocalProvider{
		userRepo: userRepo,
	}
}

func (p *LocalProvider) Authenticate(ctx context.Context, username, pwd string) (*domain.User, error) {
	// For local provider, username is email
	user, err := p.userRepo.GetByEmail(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := password.VerifyPassword(user.PasswordHash, pwd); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if user.Status != domain.UserStatusActive {
		return nil, fmt.Errorf("account is not active")
	}

	return user, nil
}

func (p *LocalProvider) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return p.userRepo.GetByID(id)
}

func (p *LocalProvider) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return p.userRepo.GetByEmail(email)
}

func (p *LocalProvider) CreateUser(ctx context.Context, user *domain.User) error {
	// Hash password if not already hashed
	if user.PasswordHash == "" {
		hashedPassword, err := password.HashPassword("") // Empty password for now
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	}

	// Set default status if not set
	if user.Status == "" {
		user.Status = domain.UserStatusPending
	}

	return p.userRepo.Create(user)
}

func (p *LocalProvider) UpdateUser(ctx context.Context, user *domain.User) error {
	return p.userRepo.Update(user)
}

func (p *LocalProvider) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return p.userRepo.Delete(id)
}

func (p *LocalProvider) ChangePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Validate password strength
	if !password.IsValidPassword(newPassword) {
		return fmt.Errorf("password does not meet requirements")
	}

	// Hash new password
	hashedPassword, err := password.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in database
	return p.userRepo.UpdatePassword(userID, hashedPassword)
}

func (p *LocalProvider) VerifyPassword(ctx context.Context, userID uuid.UUID, pwd string) (bool, error) {
	user, err := p.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	err = password.VerifyPassword(user.PasswordHash, pwd)
	return err == nil, nil
}

func (p *LocalProvider) GetProviderName() string {
	return "local"
}

func (p *LocalProvider) IsEnabled() bool {
	return true
}
