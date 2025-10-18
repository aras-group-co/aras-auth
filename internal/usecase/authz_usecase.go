package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aras-services/aras-auth/internal/domain"
)

type AuthzUseCase struct {
	roleRepo       domain.RoleRepository
	permissionRepo domain.PermissionRepository
}

func NewAuthzUseCase(roleRepo domain.RoleRepository, permissionRepo domain.PermissionRepository) *AuthzUseCase {
	return &AuthzUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

type ListRolesResponse struct {
	Roles []*domain.Role `json:"roles"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type ListPermissionsResponse struct {
	Permissions []*domain.Permission `json:"permissions"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
}

type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// Role management
func (uc *AuthzUseCase) CreateRole(ctx context.Context, req *domain.CreateRoleRequest) (*domain.Role, error) {
	role := &domain.Role{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.roleRepo.Create(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

func (uc *AuthzUseCase) GetRole(ctx context.Context, roleID uuid.UUID) (*domain.Role, error) {
	return uc.roleRepo.GetByID(roleID)
}

func (uc *AuthzUseCase) ListRoles(ctx context.Context, page, limit int) (*ListRolesResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	roles, err := uc.roleRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	total, err := uc.roleRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count roles: %w", err)
	}

	return &ListRolesResponse{
		Roles: roles,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (uc *AuthzUseCase) UpdateRole(ctx context.Context, roleID uuid.UUID, req *domain.UpdateRoleRequest) (*domain.Role, error) {
	// Get existing role
	role, err := uc.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}

	// Save updated role
	if err := uc.roleRepo.Update(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

func (uc *AuthzUseCase) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	// Check if role exists
	_, err := uc.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	// Delete role
	return uc.roleRepo.Delete(roleID)
}

func (uc *AuthzUseCase) AssignRoleToUser(ctx context.Context, userID uuid.UUID, req *domain.AssignRoleRequest) error {
	return uc.roleRepo.AssignToUser(userID, req.RoleID)
}

func (uc *AuthzUseCase) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return uc.roleRepo.RemoveFromUser(userID, roleID)
}

func (uc *AuthzUseCase) AssignRoleToGroup(ctx context.Context, groupID uuid.UUID, req *domain.AssignRoleRequest) error {
	return uc.roleRepo.AssignToGroup(groupID, req.RoleID)
}

func (uc *AuthzUseCase) RemoveRoleFromGroup(ctx context.Context, groupID, roleID uuid.UUID) error {
	return uc.roleRepo.RemoveFromGroup(groupID, roleID)
}

func (uc *AuthzUseCase) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*domain.Role, error) {
	return uc.roleRepo.GetUserRoles(userID)
}

func (uc *AuthzUseCase) GetGroupRoles(ctx context.Context, groupID uuid.UUID) ([]*domain.Role, error) {
	return uc.roleRepo.GetGroupRoles(groupID)
}

// Permission management
func (uc *AuthzUseCase) CreatePermission(ctx context.Context, req *domain.CreatePermissionRequest) (*domain.Permission, error) {
	permission := &domain.Permission{
		ID:          uuid.New(),
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.permissionRepo.Create(permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

func (uc *AuthzUseCase) GetPermission(ctx context.Context, permissionID uuid.UUID) (*domain.Permission, error) {
	return uc.permissionRepo.GetByID(permissionID)
}

func (uc *AuthzUseCase) ListPermissions(ctx context.Context, page, limit int) (*ListPermissionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	permissions, err := uc.permissionRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	total, err := uc.permissionRepo.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count permissions: %w", err)
	}

	return &ListPermissionsResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		Limit:       limit,
	}, nil
}

func (uc *AuthzUseCase) UpdatePermission(ctx context.Context, permissionID uuid.UUID, req *domain.UpdatePermissionRequest) (*domain.Permission, error) {
	// Get existing permission
	permission, err := uc.permissionRepo.GetByID(permissionID)
	if err != nil {
		return nil, fmt.Errorf("permission not found: %w", err)
	}

	// Update fields if provided
	if req.Resource != nil {
		permission.Resource = *req.Resource
	}
	if req.Action != nil {
		permission.Action = *req.Action
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}

	// Save updated permission
	if err := uc.permissionRepo.Update(permission); err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return permission, nil
}

func (uc *AuthzUseCase) DeletePermission(ctx context.Context, permissionID uuid.UUID) error {
	// Check if permission exists
	_, err := uc.permissionRepo.GetByID(permissionID)
	if err != nil {
		return fmt.Errorf("permission not found: %w", err)
	}

	// Delete permission
	return uc.permissionRepo.Delete(permissionID)
}

func (uc *AuthzUseCase) AssignPermissionToRole(ctx context.Context, roleID uuid.UUID, req *domain.AssignPermissionRequest) error {
	return uc.permissionRepo.AssignToRole(roleID, req.PermissionID)
}

func (uc *AuthzUseCase) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return uc.permissionRepo.RemoveFromRole(roleID, permissionID)
}

func (uc *AuthzUseCase) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*domain.Permission, error) {
	return uc.permissionRepo.GetRolePermissions(roleID)
}

// Authorization checks
func (uc *AuthzUseCase) CheckPermission(ctx context.Context, req *domain.CheckPermissionRequest) (*CheckPermissionResponse, error) {
	hasPermission, err := uc.permissionRepo.CheckUserPermission(req.UserID, req.Resource, req.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}

	return &CheckPermissionResponse{
		HasPermission: hasPermission,
	}, nil
}
