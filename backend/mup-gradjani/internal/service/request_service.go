package service

import (
	"errors"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("request not found")
	ErrInvalidStatus = errors.New("invalid status")
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
		Status:      domain.RequestSubmitted,
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
	case "IN_PROGRESS":
		return domain.RequestInProcess
	default:
		return domain.RequestStatus(status)
	}
}

func (s *RequestService) UpdateStatus(id string, status string) (*time.Time, error) {
	newStatus := normalizeStatus(status)

	// dozvoljavamo ova 3 kroz PATCH:
	// IN_PROCESS (ili alias IN_PROGRESS), APPROVED, REJECTED
	switch newStatus {
	case domain.RequestInProcess, domain.RequestApproved, domain.RequestRejected:
		// ok
	default:
		return nil, ErrInvalidStatus
	}

	err := s.repo.UpdateStatus(id, newStatus)
	if err != nil {
		return nil, ErrNotFound
	}

	// vraćamo processed_at samo kad je finalno stanje
	if newStatus == domain.RequestApproved || newStatus == domain.RequestRejected {
		now := time.Now().UTC()
		return &now, nil
	}

	return nil, nil
}
