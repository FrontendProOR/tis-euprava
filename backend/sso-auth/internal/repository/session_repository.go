package repository

import "tis-euprava/sso-auth/internal/domain"

type SessionRepository interface {
	FindByToken(token string) (*domain.Session, error)
	Create(session *domain.Session) error
	Delete(token string) error
}