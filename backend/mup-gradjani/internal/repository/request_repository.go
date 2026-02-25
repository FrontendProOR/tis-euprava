package repository

import (
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
)

// RequestRepository definiše ugovor ka bazi
type RequestRepository interface {
	Create(req *domain.ServiceRequest) error
	GetAll() ([]domain.ServiceRequest, error)
	FindByID(id string) (*domain.ServiceRequest, error)
	FindByCitizenID(citizenID string) ([]domain.ServiceRequest, error)
	// UpdateStatus menja status zahteva i (po potrebi) processed_at.
	// processedAt treba da bude nil za ne-finalna stanja (SUBMITTED/IN_PROCESS).
	UpdateStatus(id string, status domain.RequestStatus, processedAt *time.Time) error
}
