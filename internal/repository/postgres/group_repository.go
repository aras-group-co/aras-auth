package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aras-services/aras-auth/internal/domain"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) domain.GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *domain.Group) error {
	query := `
		INSERT INTO groups (id, name, description)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(context.Background(), query, group.ID, group.Name, group.Description)
	return err
}

func (r *GroupRepository) GetByID(id uuid.UUID) (*domain.Group, error) {
	query := `
		SELECT id, name, description, is_active, deleted_at, deleted_by, created_at, updated_at
		FROM groups WHERE id = $1 AND deleted_at IS NULL
	`

	var group domain.Group
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&group.ID, &group.Name, &group.Description, &group.IsActive, &group.DeletedAt, &group.DeletedBy, &group.CreatedAt, &group.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepository) Update(group *domain.Group) error {
	query := `
		UPDATE groups 
		SET name = $2, description = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query, group.ID, group.Name, group.Description)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

func (r *GroupRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE groups 
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("group not found or already deleted")
	}

	return nil
}

func (r *GroupRepository) List(limit, offset int) ([]*domain.Group, error) {
	query := `
		SELECT id, name, description, is_active, deleted_at, deleted_by, created_at, updated_at
		FROM groups 
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.IsActive, &group.DeletedAt, &group.DeletedBy, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}

func (r *GroupRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM groups WHERE deleted_at IS NULL`

	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}

func (r *GroupRepository) AddMember(groupID, userID uuid.UUID) error {
	query := `
		INSERT INTO user_groups (user_id, group_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, group_id) DO NOTHING
	`

	_, err := r.db.Exec(context.Background(), query, userID, groupID)
	return err
}

func (r *GroupRepository) RemoveMember(groupID, userID uuid.UUID) error {
	query := `DELETE FROM user_groups WHERE user_id = $1 AND group_id = $2`

	result, err := r.db.Exec(context.Background(), query, userID, groupID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found in group")
	}

	return nil
}

func (r *GroupRepository) GetMembers(groupID uuid.UUID) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.status, u.email_verified, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_groups ug ON u.id = ug.user_id
		WHERE ug.group_id = $1
		ORDER BY u.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.Status, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *GroupRepository) GetUserGroups(userID uuid.UUID) ([]*domain.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.created_at, g.updated_at
		FROM groups g
		INNER JOIN user_groups ug ON g.id = ug.group_id
		WHERE ug.user_id = $1
		ORDER BY g.created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	return groups, nil
}
