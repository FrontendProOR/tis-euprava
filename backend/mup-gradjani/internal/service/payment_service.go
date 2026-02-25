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

// Pay records payment and automatically advances request status (demo-friendly).
// Flow:
// - payment is stored as PAID
// - request goes PENDING -> IN_REVIEW -> (APPROVED | REJECTED)
// - APPROVED if amount >= required price for request type, otherwise REJECTED
func (s *PaymentService) Pay(requestID string, amount float64, reference string) (*domain.Payment, error) {
	if err := s.validate.ValidatePayment(requestID, amount); err != nil {
		return nil, err
	}

	// ensure request exists
	req, err := s.requests.FindByID(requestID)
	if err != nil {
		return nil, errors.New("request not found")
	}

	required, err := RequiredAmount(req.Type)
	if err != nil {
		return nil, err
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

	// store payment first
	if err := s.payments.Create(p); err != nil {
		return nil, err
	}

	// advance request status automatically
	// PENDING -> IN_REVIEW
	_ = s.requests.UpdateStatus(requestID, domain.RequestInReview, nil)

	// IN_PROCESS -> APPROVED/REJECTED
	finalStatus := domain.RequestRejected
	if amount >= required {
		finalStatus = domain.RequestApproved
	}
	processedAt := now
	if err := s.requests.UpdateStatus(requestID, finalStatus, &processedAt); err != nil {
		// payment already recorded; still return payment, but surface the error
		return p, err
	}

	return p, nil
}

func (s *PaymentService) GetByRequestID(requestID string) (*domain.Payment, error) {
	return s.payments.FindByRequestID(requestID)
}
