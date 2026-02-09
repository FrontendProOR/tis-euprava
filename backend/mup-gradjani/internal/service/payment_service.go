package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"
)

type PaymentService struct {
	payments repository.PaymentRepository
	requests repository.RequestRepository
	validate *ValidationService
}

func NewPaymentService(payRepo repository.PaymentRepository, reqRepo repository.RequestRepository) *PaymentService {
	return &PaymentService{payments: payRepo, requests: reqRepo, validate: NewValidationService()}
}

// Pay marks payment as PAID immediately (demo-friendly).
func (s *PaymentService) Pay(requestID string, amount float64, reference string) (*domain.Payment, error) {
	if err := s.validate.ValidatePayment(requestID, amount); err != nil {
		return nil, err
	}

	// ensure request exists
	if _, err := s.requests.FindByID(requestID); err != nil {
		return nil, errors.New("request not found")
	}

	now := time.Now().UTC()
	p := &domain.Payment{
		ID:        uuid.NewString(),
		RequestID: requestID,
		Amount:    amount,
		Reference: reference,
		Status:    domain.PaymentPaid,
		PaidAt:    &now,
	}

	if err := s.payments.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PaymentService) GetByRequestID(requestID string) (*domain.Payment, error) {
	return s.payments.FindByRequestID(requestID)
}
