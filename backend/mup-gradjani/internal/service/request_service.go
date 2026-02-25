package service

import (
	"errors"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("request not found")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type RequestService struct {
	repo repository.RequestRepository
}

func NewRequestService(repo repository.RequestRepository) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) CreateRequest(citizenID string, reqType string) (*domain.ServiceRequest, error) {
	req := &domain.ServiceRequest{
		ID:          uuid.NewString(),
		CitizenID:   citizenID,
		Type:        reqType,
		Status:      domain.RequestPending,
		SubmittedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *RequestService) GetAll() ([]domain.ServiceRequest, error) {
	return s.repo.GetAll()
}

func (s *RequestService) GetAllByCitizen(citizenID string) ([]domain.ServiceRequest, error) {
	return s.repo.FindByCitizenID(citizenID)
}

func (s *RequestService) GetByID(id string) (*domain.ServiceRequest, error) {
	req, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return req, nil
}

// normalizeStatus omogućava da Postman šalje IN_PROGRESS,
// a da sistem čuva IN_PROCESS (kako je u domenu)
func normalizeStatus(status string) domain.RequestStatus {
	switch status {
	case "IN_PROCESS", "IN_PROGRESS", "INREVIEW", "IN-REVIEW":
		return domain.RequestInReview
	case "PENDING":
		return domain.RequestPending
	default:
		return domain.RequestStatus(status)
	}
}

func isAllowedTransition(from, to domain.RequestStatus) bool {
	switch from {
	case domain.RequestPending:
		return to == domain.RequestInReview
	case domain.RequestInReview:
		return to == domain.RequestApproved || to == domain.RequestRejected
	case domain.RequestApproved:
		return to == domain.RequestCompleted
	case domain.RequestRejected, domain.RequestCompleted:
		return false
	default:
		return false
	}
}

// UpdateStatus menja status zahteva uz validaciju tranzicije (po dijagramu):
// PENDING -> IN_REVIEW -> (APPROVED | REJECTED) -> COMPLETED
// processed_at se setuje kada zahtev postane COMPLETED (završeno).
func (s *RequestService) UpdateStatus(id string, status string) (*domain.ServiceRequest, error) {
	newStatus := normalizeStatus(status)

	// dozvoljavamo samo statuse iz UML (po dijagramu)
	switch newStatus {
	case domain.RequestPending, domain.RequestInReview, domain.RequestApproved, domain.RequestRejected, domain.RequestCompleted:
		// ok
	default:
		return nil, ErrInvalidStatus
	}

	// učitaj trenutni zahtev da bismo proverili tranziciju
	current, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	if !isAllowedTransition(current.Status, newStatus) {
		return nil, ErrInvalidTransition
	}

	var processedAt *time.Time
	if newStatus == domain.RequestCompleted {
		now := time.Now().UTC()
		processedAt = &now
	}

	if err := s.repo.UpdateStatus(id, newStatus, processedAt); err != nil {
		return nil, ErrNotFound
	}

	current.Status = newStatus
	current.ProcessedAt = processedAt
	return current, nil
}
