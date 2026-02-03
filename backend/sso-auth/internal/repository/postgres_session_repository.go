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
	query := `
		INSERT INTO sessions
		(id, user_id, token, expires_at, created_at)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		query,
		session.ID,
		session.UserID,
		// session.Token,
		session.ExpiresAt,
		// session.CreatedAt,
	)

	return err
}

func (r *PostgresSessionRepository) FindByToken(token string) (*domain.Session, error) {
	var s domain.Session

	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = $1
	`

	err := r.db.QueryRow(query, token).Scan(
		&s.ID,
		&s.UserID,
		// &s.Token,
		&s.ExpiresAt,
		// &s.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *PostgresSessionRepository) Delete(token string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE token = $1`, token)
	return err
}

func (r *PostgresSessionRepository) DeleteExpired() error {
	_, err := r.db.Exec(
		`DELETE FROM sessions WHERE expires_at < $1`,
		time.Now(),
	)
	return err
}
