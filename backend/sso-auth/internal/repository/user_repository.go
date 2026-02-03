package repository

import "tis-euprava/sso-auth/internal/domain"

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}
