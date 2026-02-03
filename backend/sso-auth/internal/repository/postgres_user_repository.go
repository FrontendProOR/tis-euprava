package repository

import (
	"database/sql"

	"tis-euprava/sso-auth/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByUsername(username string) (*domain.User, error) {
	var u domain.User

	query := `
		SELECT id, username, password_hash, role, active, created_at
		FROM users WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
