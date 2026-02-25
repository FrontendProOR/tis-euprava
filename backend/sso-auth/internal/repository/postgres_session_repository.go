package repository

import (
	"database/sql"
	"time"

	"tis-euprava/sso-auth/internal/domain"
)

type PostgresSessionRepository struct {
	db *sql.DB
}

func NewPostgresSessionRepository(db *sql.DB) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) Create(session *domain.Session) error {
	_, err := r.db.Exec(`
		INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at)
		VALUES ($1,$2,$3,$4)
	`, session.ID, session.UserID, session.RefreshTokenHash, session.ExpiresAt)
	return err
}

func (r *PostgresSessionRepository) FindByRefreshHash(refreshHash string) (*domain.Session, error) {
	var s domain.Session
	err := r.db.QueryRow(`
		SELECT id, user_id, refresh_token_hash, expires_at, created_at, revoked_at
		FROM sessions
		WHERE refresh_token_hash = $1
	`, refreshHash).Scan(
		&s.ID, &s.UserID, &s.RefreshTokenHash, &s.ExpiresAt, &s.CreatedAt, &s.RevokedAt,
	)
	if err != nil {
		return nil, err
	}

	if s.RevokedAt != nil || time.Now().After(s.ExpiresAt) {
		return nil, sql.ErrNoRows
	}

	return &s, nil
}

func (r *PostgresSessionRepository) RevokeByID(id string) error {
	_, err := r.db.Exec(`UPDATE sessions SET revoked_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *PostgresSessionRepository) RevokeByRefreshHash(refreshHash string) error {
	_, err := r.db.Exec(`UPDATE sessions SET revoked_at=NOW() WHERE refresh_token_hash=$1`, refreshHash)
	return err
}

func (r *PostgresSessionRepository) DeleteExpired() error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE expires_at < $1`, time.Now())
	return err
}