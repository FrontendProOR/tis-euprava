package repository

import "tis-euprava/sso-auth/internal/domain"

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	Create(user *domain.User) error
	ExistsByUsername(username string) (bool, error)
}