package domain

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required,min=1,max=100"`
	Description string     `json:"description" db:"description"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy   *uuid.UUID `json:"deleted_by,omitempty" db:"deleted_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateGroupRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty"`
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type GroupRepository interface {
	Create(group *Group) error
	GetByID(id uuid.UUID) (*Group, error)
	Update(group *Group) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Group, error)
	Count() (int, error)
	AddMember(groupID, userID uuid.UUID) error
	RemoveMember(groupID, userID uuid.UUID) error
	GetMembers(groupID uuid.UUID) ([]*User, error)
	GetUserGroups(userID uuid.UUID) ([]*Group, error)
}
