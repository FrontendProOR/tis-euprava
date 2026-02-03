package repository

import (
	"database/sql"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
)

type PostgresRequestRepository struct {
	db *sql.DB
}

func NewPostgresRequestRepository(db *sql.DB) *PostgresRequestRepository {
	return &PostgresRequestRepository{db: db}
}

func (r *PostgresRequestRepository) Create(req *domain.ServiceRequest) error {
	query := `
		INSERT INTO service_requests
		(id, citizen_id, type, status, submitted_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.Exec(
		query,
		req.ID,
		req.CitizenID,
		req.Type,
		req.Status,
		req.SubmittedAt,
		// req.UpdatedAt,
	)

	return err
}

func (r *PostgresRequestRepository) FindByID(id string) (*domain.ServiceRequest, error) {
	var req domain.ServiceRequest

	query := `
		SELECT id, citizen_id, type, status, submitted_at, updated_at
		FROM service_requests
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&req.CitizenID,
		&req.Type,
		&req.Status,
		&req.SubmittedAt,
		// &req.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *PostgresRequestRepository) FindByCitizenID(citizenID string) ([]domain.ServiceRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, citizen_id, type, status, submitted_at, updated_at
		FROM service_requests
		WHERE citizen_id = $1
	`, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.ServiceRequest

	for rows.Next() {
		var req domain.ServiceRequest
		if err := rows.Scan(
			&req.ID,
			&req.CitizenID,
			&req.Type,
			&req.Status,
			&req.SubmittedAt,
			// &req.UpdatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *PostgresRequestRepository) UpdateStatus(id string, status domain.RequestStatus) error {
	query := `
		UPDATE service_requests
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}
