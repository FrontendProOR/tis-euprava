package repository

import "tis-euprava/zzs-statistika/internal/domain"

type SurveyRepository interface {
	Create(survey *domain.Survey) error
	FindAll() ([]domain.Survey, error)
	//FindByID(id string) (*domain.Survey, error)
}
