package repository

import "tis-euprava/mup-gradjani/internal/domain"

// RequestRepository definiše ugovor ka bazi
type RequestRepository interface {
	Create(req *domain.ServiceRequest) error
	GetAll() ([]domain.ServiceRequest, error)
	FindByID(id string) (*domain.ServiceRequest, error)
	FindByCitizenID(citizenID string) ([]domain.ServiceRequest, error)
	UpdateStatus(id string, status domain.RequestStatus) error
}
