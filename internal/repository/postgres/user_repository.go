package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aras-services/aras-auth/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, status, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(context.Background(), query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Status, user.EmailVerified)

	return err
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, status, email_verified, is_deleted, is_system, created_at, updated_at
		FROM users WHERE id = $1 AND is_deleted = FALSE
	`

	var user domain.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Status, &user.EmailVerified, &user.IsDeleted, &user.IsSystem, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, status, email_verified, is_deleted, is_system, created_at, updated_at
		FROM users WHERE email = $1 AND is_deleted = FALSE
	`

	var user domain.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Status, &user.EmailVerified, &user.IsDeleted, &user.IsSystem, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $2, first_name = $3, last_name = $4, status = $5, email_verified = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query,
		user.ID, user.Email, user.FirstName, user.LastName, user.Status, user.EmailVerified)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `
		UPDATE users 
		SET is_deleted = TRUE, updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

func (r *UserRepository) List(limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, status, email_verified, is_deleted, is_system, created_at, updated_at
		FROM users 
		WHERE is_deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(context.Background(), query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
			&user.Status, &user.EmailVerified, &user.IsDeleted, &user.IsSystem, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_deleted = FALSE`

	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}

func (r *UserRepository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id, passwordHash)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) UpdateEmailVerified(id uuid.UUID, verified bool) error {
	query := `UPDATE users SET email_verified = $2, updated_at = NOW() WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id, verified)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
