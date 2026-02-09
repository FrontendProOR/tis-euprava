package repository

import (
	"database/sql"
	"errors"

	"tis-euprava/mup-gradjani/internal/domain"
)

type PostgresCertificateRepository struct {
	db *sql.DB
}

func NewPostgresCertificateRepository(db *sql.DB) *PostgresCertificateRepository {
	return &PostgresCertificateRepository{db: db}
}

// compile-time check
var _ CertificateRepository = (*PostgresCertificateRepository)(nil)

func (r *PostgresCertificateRepository) Create(c *domain.ElectronicCertificate) error {
	if c == nil {
		return errors.New("certificate is nil")
	}
	if c.ID == "" || c.RequestID == "" {
		return errors.New("certificate id/requestId is empty")
	}
	if len(c.PDF) == 0 {
		return errors.New("certificate pdf is empty")
	}

	_, err := r.db.Exec(`
		INSERT INTO certificates (id, request_id, issued_at, pdf)
		VALUES ($1,$2,$3,$4)
	`, c.ID, c.RequestID, c.IssuedAt, c.PDF)

	return err
}

func (r *PostgresCertificateRepository) FindByRequestID(requestID string) (*domain.ElectronicCertificate, error) {
	if requestID == "" {
		return nil, errors.New("requestID is empty")
	}

	var c domain.ElectronicCertificate
	err := r.db.QueryRow(`
		SELECT id, request_id, issued_at, pdf
		FROM certificates
		WHERE request_id = $1
	`, requestID).Scan(&c.ID, &c.RequestID, &c.IssuedAt, &c.PDF)

	if err != nil {
		// bitno: prosledi sql.ErrNoRows gore (service ga često posebno tretira)
		return nil, err
	}

	return &c, nil
}
