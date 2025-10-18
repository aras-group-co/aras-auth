package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aras-services/aras-auth/internal/domain"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) domain.RefreshTokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query, token.ID, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (r *TokenRepository) GetByID(id uuid.UUID) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE id = $1
	`

	var token domain.RefreshToken
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) GetByUserID(userID uuid.UUID) ([]*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens 
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*domain.RefreshToken
	for rows.Next() {
		var token domain.RefreshToken
		err := rows.Scan(
			&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

func (r *TokenRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("refresh token not found")
	}

	return nil
}

func (r *TokenRepository) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := r.db.Exec(context.Background(), query, userID)
	return err
}

func (r *TokenRepository) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`

	_, err := r.db.Exec(context.Background(), query)
	return err
}

func (r *TokenRepository) GetByTokenHash(tokenHash string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens 
		WHERE token_hash = $1 AND expires_at > NOW()
	`

	var token domain.RefreshToken
	err := r.db.QueryRow(context.Background(), query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) CleanupExpiredTokens() (int, error) {
	query := `SELECT cleanup_expired_tokens()`

	var count int
	err := r.db.QueryRow(context.Background(), query).Scan(&count)
	return count, err
}


