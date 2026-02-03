package	repository

import (
	// "database/sql"
	"tis-euprava/zzs-statistika/internal/domain"
)
type ViolationRepository interface {
	Create(v *domain.Violation) error
	FindAll() ([]domain.Violation, error)
	//FindByResidence(residenceID string) ([]domain.Violation, error)
}