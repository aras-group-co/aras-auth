package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aras-services/aras-auth/internal/domain"
)

type PermissionRepository struct {
	db *pgxpool.Pool
}

func NewPermissionRepository(db *pgxpool.Pool) domain.PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (id, resource, action, description)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query, permission.ID, permission.Resource, permission.Action, permission.Description)
	return err
}

func (r *PermissionRepository) GetByID(id uuid.UUID) (*domain.Permission, error) {
	query := `
		SELECT id, resource, action, description, is_active, is_deleted, is_system, created_at, updated_at
		FROM permissions WHERE id = $1 AND is_deleted = FALSE
	`

	var permission domain.Permission
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&permission.ID, &permission.Resource, &permission.Action, &permission.Description, &permission.IsActive, &permission.IsDeleted, &permission.IsSystem, &permission.CreatedAt, &permission.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) GetByResourceAndAction(resource, action string) (*domain.Permission, error) {
	query := `
		SELECT id, resource, action, description, is_active, is_deleted, is_system, created_at, updated_at
		FROM permissions WHERE resource = $1 AND action = $2 AND is_deleted = FALSE
	`

	var permission domain.Permission
	err := r.db.QueryRow(context.Background(), query, resource, action).Scan(
		&permission.ID, &permission.Resource, &permission.Action, &permission.Description, &permission.IsActive, &permission.IsDeleted, &permission.IsSystem, &permission.CreatedAt, &permission.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

func (r *PermissionRepository) Update(permission *domain.Permission) error {
	query := `
		UPDATE permissions 
		SET resource = $2, action = $3, description = $4, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query, permission.ID, permission.Resource, permission.Action, permission.Description)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("permission not found")
	}

	return nil
}

func (r *PermissionRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE permissions 
		SET is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("permission not found or already deleted")
	}

	return nil
}

func (r *PermissionRepository) List(limit, offset int) ([]*domain.Permission, error) {
	query := `
		SELECT id, resource, action, description, is_active, is_deleted, is_system, created_at, updated_at
		FROM permissions 
		WHERE is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.ID, &permission.Resource, &permission.Action, &permission.Description, &permission.IsActive, &permission.IsDeleted, &permission.IsSystem, &permission.CreatedAt, &permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

func (r *PermissionRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM permissions WHERE is_deleted = FALSE`

	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}

func (r *PermissionRepository) AssignToRole(roleID, permissionID uuid.UUID) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.Exec(context.Background(), query, roleID, permissionID)
	return err
}

func (r *PermissionRepository) RemoveFromRole(roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.db.Exec(context.Background(), query, roleID, permissionID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("permission not assigned to role")
	}

	return nil
}

func (r *PermissionRepository) GetRolePermissions(roleID uuid.UUID) ([]*domain.Permission, error) {
	query := `
		SELECT p.id, p.resource, p.action, p.description, p.is_active, p.is_deleted, p.is_system, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		  AND p.is_deleted = FALSE
		  AND p.is_active = TRUE
		ORDER BY p.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.ID, &permission.Resource, &permission.Action, &permission.Description, &permission.IsActive, &permission.IsDeleted, &permission.IsSystem, &permission.CreatedAt, &permission.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, nil
}

func (r *PermissionRepository) CheckUserPermission(userID uuid.UUID, resource, action string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 
		  AND p.resource = $2 
		  AND p.action = $3
		  AND p.is_deleted = FALSE
		  AND p.is_active = TRUE
		  AND r.is_deleted = FALSE
		  AND r.is_active = TRUE
		
		UNION ALL
		
		SELECT COUNT(*) > 0
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN group_roles gr ON r.id = gr.role_id
		INNER JOIN groups g ON gr.group_id = g.id
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1 
		  AND p.resource = $2 
		  AND p.action = $3
		  AND p.is_deleted = FALSE
		  AND p.is_active = TRUE
		  AND r.is_deleted = FALSE
		  AND r.is_active = TRUE
		  AND g.is_deleted = FALSE
		  AND g.is_active = TRUE
	`

	rows, err := r.db.Query(context.Background(), query, userID, resource, action)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var hasPermission bool
		err := rows.Scan(&hasPermission)
		if err != nil {
			return false, err
		}
		if hasPermission {
			return true, nil
		}
	}

	return false, nil
}
