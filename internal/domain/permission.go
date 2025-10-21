package domain

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Resource    string     `json:"resource" db:"resource" validate:"required,min=1,max=100"`
	Action      string     `json:"action" db:"action" validate:"required,min=1,max=100"`
	Description string     `json:"description" db:"description"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreatePermissionRequest struct {
	Resource    string `json:"resource" validate:"required,min=1,max=100"`
	Action      string `json:"action" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Resource    *string `json:"resource,omitempty" validate:"omitempty,min=1,max=100"`
	Action      *string `json:"action,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
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
