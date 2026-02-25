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

func (r *PostgresUserRepository) ExistsByUsername(username string) (bool, error) {
	row := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`, username)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgresUserRepository) Create(user *domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (id, username, password_hash, role, active)
		VALUES ($1,$2,$3,$4,$5)
	`, user.ID, user.Username, user.PasswordHash, string(user.Role), user.Active)
	return err
}

func (r *PostgresUserRepository) FindByUsername(username string) (*domain.User, error) {
	var u domain.User
	query := `
		SELECT id, username, password_hash, role, active, created_at
		FROM users WHERE username = $1
	`
	err := r.db.QueryRow(query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PostgresUserRepository) FindByID(id string) (*domain.User, error) {
	var u domain.User
	query := `
		SELECT id, username, password_hash, role, active, created_at
		FROM users WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Active, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}