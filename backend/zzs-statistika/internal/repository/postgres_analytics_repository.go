package repository

import (
	"database/sql"

	"tis-euprava/zzs-statistika/internal/domain"
)

type PostgresAnalyticsRepository struct {
	db *sql.DB
}

func NewPostgresAnalyticsRepository(db *sql.DB) *PostgresAnalyticsRepository {
	return &PostgresAnalyticsRepository{db: db}
}

func (r *PostgresAnalyticsRepository) GetViolationCountByType() ([]domain.AnalyticsResult, error) {
	rows, err := r.db.Query(`
		SELECT type, COUNT(*) as count
		FROM violations
		GROUP BY type
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AnalyticsResult

	for rows.Next() {
		var res domain.AnalyticsResult
		if err := rows.Scan(&res.Type, &res.Count); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}

func (r *PostgresAnalyticsRepository) GetSurveyParticipation() ([]domain.AnalyticsResult, error) {
	rows, err := r.db.Query(`
		SELECT survey_id, COUNT(*) as count
		FROM survey_responses
		GROUP BY survey_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.AnalyticsResult

	for rows.Next() {
		var res domain.AnalyticsResult
		if err := rows.Scan(&res.SurveyID, &res.Count); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}
