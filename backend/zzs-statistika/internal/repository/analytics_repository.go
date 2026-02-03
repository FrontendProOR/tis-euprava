package repository

import "tis-euprava/zzs-statistika/internal/domain"

type AnalyticsRepository interface {
	GetViolationCountByType() ([]domain.AnalyticsResult, error)
	GetSurveyParticipation() ([]domain.AnalyticsResult, error)
}
