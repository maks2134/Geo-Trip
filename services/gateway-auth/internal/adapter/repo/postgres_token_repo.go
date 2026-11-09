package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/makskozlov/atlasflow/services/gateway-auth/internal/domain"
)

type PostgresTokenRepository struct {
	db *sql.DB
}

func NewPostgresTokenRepository(db *sql.DB) *PostgresTokenRepository {
	return &PostgresTokenRepository{db: db}
}

func (r *PostgresTokenRepository) Save(ctx context.Context, record *domain.RefreshTokenRecord) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, revoked_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		record.ID,
		record.UserID,
		record.TokenHash,
		record.ExpiresAt,
		record.RevokedAt,
		record.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

func (r *PostgresTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenRecord, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var rec domain.RefreshTokenRecord
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.TokenHash,
		&rec.ExpiresAt,
		&rec.RevokedAt,
		&rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrInvalidToken
		}
		return nil, fmt.Errorf("query refresh token: %w", err)
	}
	return &rec, nil
}

func (r *PostgresTokenRepository) Revoke(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token_hash = $2 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, revokedAt, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

func (r *PostgresTokenRepository) RevokeAllForUser(ctx context.Context, userID string, revokedAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE user_id = $2 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, revokedAt, userID)
	if err != nil {
		return fmt.Errorf("revoke user refresh tokens: %w", err)
	}
	return nil
}