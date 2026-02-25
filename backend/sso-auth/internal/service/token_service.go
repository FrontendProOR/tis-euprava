package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"tis-euprava/sso-auth/internal/config"
	"tis-euprava/sso-auth/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type TokenService struct {
	cfg *config.Config
}

func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{cfg: cfg}
}

func (t *TokenService) NewRefreshToken() string {
	return uuid.NewString() + "-" + uuid.NewString()
}

func (t *TokenService) HashRefreshToken(refresh string) string {
	sum := sha256.Sum256([]byte(refresh))
	return hex.EncodeToString(sum[:])
}

func (t *TokenService) SignAccessToken(u *domain.User) (string, time.Time, error) {
	exp := time.Now().Add(time.Duration(t.cfg.AccessTTLMin) * time.Minute)

	claims := jwt.MapClaims{
		"iss":  t.cfg.JWTIssuer,
		"sub":  u.ID,
		"usr":  u.Username,
		"role": string(u.Role),
		"exp":  exp.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(t.cfg.JWTSecret))
	return signed, exp, err
}

func (t *TokenService) ParseAccessToken(tokenStr string) (map[string]any, error) {
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(t.cfg.JWTSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	out := map[string]any{
		"sub":  claims["sub"],
		"usr":  claims["usr"],
		"role": claims["role"],
		"iss":  claims["iss"],
	}
	return out, nil
}