package repository

import "tis-euprava/mup-gradjani/internal/domain"

type AppointmentRepository interface {
	Create(a *domain.Appointment) error
	FindByCitizenID(citizenID string) ([]domain.Appointment, error)
	Delete(id string) error
}
