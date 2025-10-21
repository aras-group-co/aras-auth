package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusPending   UserStatus = "pending"
	UserStatusSuspended UserStatus = "suspended"
)

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email" validate:"required,email"`
	PasswordHash  string     `json:"-" db:"password_hash"`
	FirstName     string     `json:"first_name" db:"first_name"`
	LastName      string     `json:"last_name" db:"last_name"`
	Status        UserStatus `json:"status" db:"status"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy     *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserRequest struct {
	FirstName *string     `json:"first_name,omitempty"`
	LastName  *string     `json:"last_name,omitempty"`
	Status    *UserStatus `json:"status,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*User, error)
	Count() (int, error)
	UpdatePassword(id uuid.UUID, passwordHash string) error
	UpdateEmailVerified(id uuid.UUID, verified bool) error
}
