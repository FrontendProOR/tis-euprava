package repository

import "tis-euprava/sso-auth/internal/domain"

type SessionRepository interface {
	Create(session *domain.Session) error
	FindByRefreshHash(refreshHash string) (*domain.Session, error)
	RevokeByID(id string) error
	RevokeByRefreshHash(refreshHash string) error
	DeleteExpired() error
}
