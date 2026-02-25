package service

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"tis-euprava/sso-auth/internal/config"
	"tis-euprava/sso-auth/internal/domain"
	"tis-euprava/sso-auth/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrBadCredentials    = errors.New("bad credentials")
	ErrUserInactive      = errors.New("user inactive")
	ErrInvalidRefresh    = errors.New("invalid refresh token")
)

type AuthService struct {
	cfg      *config.Config
	users    repository.UserRepository
	sessions repository.SessionRepository
	tokens   *TokenService
}

func NewAuthService(cfg *config.Config, users repository.UserRepository, sessions repository.SessionRepository, tokens *TokenService) *AuthService {
	return &AuthService{cfg: cfg, users: users, sessions: sessions, tokens: tokens}
}

func (s *AuthService) Register(username, password string, role domain.Role) (*domain.User, error) {
	username = strings.TrimSpace(strings.ToLower(username))
	if username == "" || password == "" {
		return nil, ErrBadCredentials
	}

	exists, err := s.users.ExistsByUsername(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		Active:       true,
	}

	if err := s.users.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(username, password, userAgent, ip string) (domain.TokenPair, *domain.User, error) {
	username = strings.TrimSpace(strings.ToLower(username))

	u, err := s.users.FindByUsername(username)
	if err != nil {
		return domain.TokenPair{}, nil, ErrBadCredentials
	}
	if !u.Active {
		return domain.TokenPair{}, nil, ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return domain.TokenPair{}, nil, ErrBadCredentials
	}

	access, exp, err := s.tokens.SignAccessToken(u)
	if err != nil {
		return domain.TokenPair{}, nil, err
	}

	refresh := s.tokens.NewRefreshToken()
	refreshHash := s.tokens.HashRefreshToken(refresh)

	session := &domain.Session{
		ID:               uuid.NewString(),
		UserID:           u.ID,
		RefreshTokenHash: refreshHash,
		ExpiresAt:        time.Now().Add(time.Duration(s.cfg.RefreshTTLDays) * 24 * time.Hour),
	}

	if err := s.sessions.Create(session); err != nil {
		return domain.TokenPair{}, nil, err
	}

	return domain.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    exp,
	}, u, nil
}

func (s *AuthService) Refresh(refreshToken string) (domain.TokenPair, error) {
	refreshHash := s.tokens.HashRefreshToken(strings.TrimSpace(refreshToken))

	sess, err := s.sessions.FindByRefreshHash(refreshHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TokenPair{}, ErrInvalidRefresh
		}
		return domain.TokenPair{}, err
	}

	u, err := s.users.FindByID(sess.UserID)
	if err != nil {
		return domain.TokenPair{}, ErrInvalidRefresh
	}
	if !u.Active {
		return domain.TokenPair{}, ErrUserInactive
	}

	access, exp, err := s.tokens.SignAccessToken(u)
	if err != nil {
		return domain.TokenPair{}, err
	}

	// rotate refresh: revoke old and create new
	_ = s.sessions.RevokeByID(sess.ID)

	newRefresh := s.tokens.NewRefreshToken()
	newHash := s.tokens.HashRefreshToken(newRefresh)

	newSess := &domain.Session{
		ID:               uuid.NewString(),
		UserID:           u.ID,
		RefreshTokenHash: newHash,
		ExpiresAt:        time.Now().Add(time.Duration(s.cfg.RefreshTTLDays) * 24 * time.Hour),
	}

	if err := s.sessions.Create(newSess); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  access,
		RefreshToken: newRefresh,
		ExpiresAt:    exp,
	}, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	refreshHash := s.tokens.HashRefreshToken(strings.TrimSpace(refreshToken))
	return s.sessions.RevokeByRefreshHash(refreshHash)
}