package repository

import "tis-euprava/mup-gradjani/internal/domain"

type PaymentRepository interface {
	Create(p *domain.Payment) error
	FindByRequestID(requestID string) (*domain.Payment, error)
}
