package repository

import "tis-euprava/mup-gradjani/internal/domain"

type CertificateRepository interface {
	Create(c *domain.ElectronicCertificate) error
	FindByRequestID(requestID string) (*domain.ElectronicCertificate, error)
}
