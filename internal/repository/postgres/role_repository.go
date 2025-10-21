package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aras-services/aras-auth/internal/domain"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(db *pgxpool.Pool) domain.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, description)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(context.Background(), query, role.ID, role.Name, role.Description)
	return err
}

func (r *RoleRepository) GetByID(id uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, is_deleted, created_at, updated_at
		FROM roles WHERE id = $1 AND is_deleted = FALSE
	`

	var role domain.Role
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.IsActive, &role.IsDeleted, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetByName(name string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, is_deleted, created_at, updated_at
		FROM roles WHERE name = $1 AND is_deleted = FALSE
	`

	var role domain.Role
	err := r.db.QueryRow(context.Background(), query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.IsActive, &role.IsDeleted, &role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) Update(role *domain.Role) error {
	query := `
		UPDATE roles 
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query, role.ID, role.Name, role.Description)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (r *RoleRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE roles 
		SET is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not found or already deleted")
	}

	return nil
}

func (r *RoleRepository) List(limit, offset int) ([]*domain.Role, error) {
	query := `
		SELECT id, name, description, is_active, is_deleted, created_at, updated_at
		FROM roles 
		WHERE is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.IsActive, &role.IsDeleted, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *RoleRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM roles WHERE is_deleted = FALSE`

	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}

func (r *RoleRepository) AssignToUser(userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.Exec(context.Background(), query, userID, roleID)
	return err
}

func (r *RoleRepository) RemoveFromUser(userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.db.Exec(context.Background(), query, userID, roleID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not assigned to user")
	}

	return nil
}

func (r *RoleRepository) AssignToGroup(groupID, roleID uuid.UUID) error {
	query := `
		INSERT INTO group_roles (group_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (group_id, role_id) DO NOTHING
	`

	_, err := r.db.Exec(context.Background(), query, groupID, roleID)
	return err
}

func (r *RoleRepository) RemoveFromGroup(groupID, roleID uuid.UUID) error {
	query := `DELETE FROM group_roles WHERE group_id = $1 AND role_id = $2`

	result, err := r.db.Exec(context.Background(), query, groupID, roleID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not assigned to group")
	}

	return nil
}

func (r *RoleRepository) GetUserRoles(userID uuid.UUID) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}

func (r *RoleRepository) GetGroupRoles(groupID uuid.UUID) ([]*domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN group_roles gr ON r.id = gr.role_id
		WHERE gr.group_id = $1
		ORDER BY r.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}
