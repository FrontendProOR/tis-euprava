package repository

import (
	"database/sql"

	"tis-euprava/zzs-statistika/internal/domain"
)

type PostgresSurveyRepository struct {
	db *sql.DB
}

func NewPostgresSurveyRepository(db *sql.DB) *PostgresSurveyRepository {
	return &PostgresSurveyRepository{db: db}
}

func (r *PostgresSurveyRepository) Create(s *domain.Survey) error {
	query := `
		INSERT INTO surveys (id, title, description, created_at, closed)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		query,
		s.ID,
		s.Title,
		s.Description,
		s.CreatedAt,
		s.Closed,
	)

	return err
}

func (r *PostgresSurveyRepository) FindAll() ([]domain.Survey, error) {
	rows, err := r.db.Query(`
		SELECT id, title, description, created_at, closed FROM surveys
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var surveys []domain.Survey
	for rows.Next() {
		var s domain.Survey
		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Description,
			&s.CreatedAt,
			&s.Closed,
		); err != nil {
			return nil, err
		}
		surveys = append(surveys, s)
	}

	return surveys, nil
}
