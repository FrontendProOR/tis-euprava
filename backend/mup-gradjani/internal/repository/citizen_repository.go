package repository

import "tis-euprava/mup-gradjani/internal/domain"

type CitizenRepository interface {
	Create(citizen *domain.Citizen) error
	FindByID(id string) (*domain.Citizen, error)
	FindByJMBG(jmbg string) (*domain.Citizen, error)
	List() ([]domain.Citizen, error)
}
