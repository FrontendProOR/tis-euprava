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

// compile-time check
var _ RequestRepository = (*PostgresRequestRepository)(nil)

func (r *PostgresRequestRepository) Create(req *domain.ServiceRequest) error {
	query := `
		INSERT INTO service_requests
		(id, citizen_id, type, status, submitted_at)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		query,
		req.ID,
		req.CitizenID,
		req.Type,
		req.Status,
		req.SubmittedAt,
	)

	return err
}

func (r *PostgresRequestRepository) GetAll() ([]domain.ServiceRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, citizen_id, type, status, submitted_at, processed_at
		FROM service_requests
		ORDER BY submitted_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.ServiceRequest
	for rows.Next() {
		var req domain.ServiceRequest
		if err := rows.Scan(
			&req.ID,
			&req.CitizenID,
			&req.Type,
			&req.Status,
			&req.SubmittedAt,
			&req.ProcessedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, req)
	}

	return list, rows.Err()
}

func (r *PostgresRequestRepository) FindByID(id string) (*domain.ServiceRequest, error) {
	var req domain.ServiceRequest

	query := `
		SELECT id, citizen_id, type, status, submitted_at, processed_at
		FROM service_requests
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&req.CitizenID,
		&req.Type,
		&req.Status,
		&req.SubmittedAt,
		&req.ProcessedAt,
	)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (r *PostgresRequestRepository) FindByCitizenID(citizenID string) ([]domain.ServiceRequest, error) {
	rows, err := r.db.Query(`
		SELECT id, citizen_id, type, status, submitted_at, processed_at
		FROM service_requests
		WHERE citizen_id = $1
	`, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.ServiceRequest
	for rows.Next() {
		var req domain.ServiceRequest
		if err := rows.Scan(
			&req.ID,
			&req.CitizenID,
			&req.Type,
			&req.Status,
			&req.SubmittedAt,
			&req.ProcessedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, req)
	}

	return list, rows.Err()
}

// UpdateStatus:
// - uvek update-uje status
// - processed_at se setuje SAMO za APPROVED/REJECTED, inače NULL
func (r *PostgresRequestRepository) UpdateStatus(id string, status domain.RequestStatus, processedAt *time.Time) error {
	var processed any = nil
	if processedAt != nil {
		t := (*processedAt).UTC()
		processed = t
	}

	query := `
		UPDATE service_requests
		SET status = $1, processed_at = $2
		WHERE id = $3
	`

	res, err := r.db.Exec(query, status, processed, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err == nil && affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
