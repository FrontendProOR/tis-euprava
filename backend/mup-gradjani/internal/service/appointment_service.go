package service

import (
	"time"

	"github.com/google/uuid"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"
)

type AppointmentService struct {
	repo repository.AppointmentRepository
}

func NewAppointmentService(repo repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{repo: repo}
}

func (s *AppointmentService) Schedule(citizenID, policeStation string, dt time.Time) (*domain.Appointment, error) {
	a := &domain.Appointment{
		ID:            uuid.NewString(),
		CitizenID:     citizenID,
		PoliceStation: policeStation,
		DateTime:      dt,
		Status:        domain.AppointmentScheduled,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AppointmentService) ListByCitizen(citizenID string) ([]domain.Appointment, error) {
	return s.repo.FindByCitizenID(citizenID)
}

func (s *AppointmentService) Cancel(id string) error { return s.repo.Delete(id) }
