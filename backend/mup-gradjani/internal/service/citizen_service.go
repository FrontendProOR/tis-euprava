package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/repository"
)

var (
	ErrCitizenNotFound  = errors.New("citizen not found")
	ErrCitizenDuplicate = errors.New("citizen already exists")
)

type CitizenService struct {
	repo     repository.CitizenRepository
	validate *ValidationService
}

func NewCitizenService(repo repository.CitizenRepository) *CitizenService {
	return &CitizenService{repo: repo, validate: NewValidationService()}
}

func (s *CitizenService) Create(c *domain.Citizen) (*domain.Citizen, error) {
	if err := s.validate.ValidateCitizen(c); err != nil {
		return nil, err
	}

	// prevent duplicates by JMBG
	existing, err := s.repo.FindByJMBG(c.JMBG)
	if err == nil && existing != nil {
		return nil, ErrCitizenDuplicate
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// real db error
		return nil, err
	}

	c.ID = uuid.NewString()
	c.CreatedAt = time.Now().UTC()

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CitizenService) GetByID(id string) (*domain.Citizen, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrCitizenNotFound
	}
	return c, nil
}

func (s *CitizenService) List() ([]domain.Citizen, error) {
	return s.repo.List()
}
