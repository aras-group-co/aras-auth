package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Description string    `json:"description" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	IsSystem    bool      `json:"is_system" db:"is_system"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}

type RoleRepository interface {
	Create(role *Role) error
	GetByID(id uuid.UUID) (*Role, error)
	GetByName(name string) (*Role, error)
	Update(role *Role) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Role, error)
	Count() (int, error)
	AssignToUser(userID, roleID uuid.UUID) error
	RemoveFromUser(userID, roleID uuid.UUID) error
	AssignToGroup(groupID, roleID uuid.UUID) error
	RemoveFromGroup(groupID, roleID uuid.UUID) error
	GetUserRoles(userID uuid.UUID) ([]*Role, error)
	GetGroupRoles(groupID uuid.UUID) ([]*Role, error)
}
