package domain

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Resource    string    `json:"resource" db:"resource" validate:"required,min=1,max=100"`
	Action      string    `json:"action" db:"action" validate:"required,min=1,max=100"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreatePermissionRequest struct {
	Resource    string `json:"resource" validate:"required,min=1,max=100"`
	Action      string `json:"action" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active,omitempty"` // optional, default true
}

type UpdatePermissionRequest struct {
	Resource    *string `json:"resource,omitempty" validate:"omitempty,min=1,max=100"`
	Action      *string `json:"action,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type AssignPermissionRequest struct {
	PermissionID uuid.UUID `json:"permission_id" validate:"required"`
}

type CheckPermissionRequest struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Resource string    `json:"resource" validate:"required"`
	Action   string    `json:"action" validate:"required"`
}

type PermissionRepository interface {
	Create(permission *Permission) error
	GetByID(id uuid.UUID) (*Permission, error)
	GetByResourceAndAction(resource, action string) (*Permission, error)
	Update(permission *Permission) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Permission, error)
	Count() (int, error)
	AssignToRole(roleID, permissionID uuid.UUID) error
	RemoveFromRole(roleID, permissionID uuid.UUID) error
	GetRolePermissions(roleID uuid.UUID) ([]*Permission, error)
	CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error)
}
