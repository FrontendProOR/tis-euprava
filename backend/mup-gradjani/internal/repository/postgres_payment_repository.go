package repository

import (
	"database/sql"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
)

type PostgresPaymentRepository struct{ db *sql.DB }

func NewPostgresPaymentRepository(db *sql.DB) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{db: db}
}

var _ PaymentRepository = (*PostgresPaymentRepository)(nil)

func (r *PostgresPaymentRepository) Create(p *domain.Payment) error {
	q := `INSERT INTO payments (id, request_id, amount, reference, status, paid_at) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.db.Exec(q, p.ID, p.RequestID, p.Amount, p.Reference, p.Status, p.PaidAt)
	return err
}

func (r *PostgresPaymentRepository) FindByRequestID(requestID string) (*domain.Payment, error) {
	var p domain.Payment
	var paidAt sql.NullTime
	q := `SELECT id, request_id, amount, reference, status, paid_at FROM payments WHERE request_id=$1 ORDER BY paid_at DESC NULLS LAST LIMIT 1`
	err := r.db.QueryRow(q, requestID).Scan(&p.ID, &p.RequestID, &p.Amount, &p.Reference, &p.Status, &paidAt)
	if err != nil {
		return nil, err
	}
	if paidAt.Valid {
		t := paidAt.Time
		p.PaidAt = &t
	}
	// normalize if missing paid_at but status paid
	if p.Status == domain.PaymentPaid && p.PaidAt == nil {
		n := time.Now().UTC()
		p.PaidAt = &n
	}
	return &p, nil
}
