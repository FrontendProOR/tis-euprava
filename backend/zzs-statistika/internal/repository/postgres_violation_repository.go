package repository

import (
	"database/sql"

	"tis-euprava/zzs-statistika/internal/domain"
)

type PostgresViolationRepository struct {
	db *sql.DB
}

func NewPostgresViolationRepository(db *sql.DB) *PostgresViolationRepository {
	return &PostgresViolationRepository{db: db}
}

func (r *PostgresViolationRepository) Create(v *domain.Violation) error {
	query := `
		INSERT INTO violations
		(id, type, residence_id, occurred_at, reported_at)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		query,
		v.ID,
		v.Type,
		// v.ResidenceID,
		v.OccurredAt,
		// v.ReportedAt,
	)

	return err
}

func (r *PostgresViolationRepository) FindAll() ([]domain.Violation, error) {
	rows, err := r.db.Query(`
		SELECT id, type, residence_id, occurred_at, reported_at
		FROM violations
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.Violation

	for rows.Next() {
		var v domain.Violation
		if err := rows.Scan(
			&v.ID,
			&v.Type,
			// &v.ResidenceID,
			&v.OccurredAt,
			// &v.ReportedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, v)
	}

	return list, nil
}
