package service

import (
	"strings"

	"tis-euprava/sso-auth/internal/domain"
)

type SSOService struct {
	tokens *TokenService
}

func NewSSOService(tokens *TokenService) *SSOService {
	return &SSOService{tokens: tokens}
}

func (s *SSOService) ValidateBearer(authHeader string) (map[string]any, error) {
	authHeader = strings.TrimSpace(authHeader)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrInvalidToken
	}
	raw := strings.TrimPrefix(authHeader, "Bearer ")
	return s.tokens.ParseAccessToken(raw)
}

// convenience (optional for future)
func (s *SSOService) RoleFromClaims(claims map[string]any) domain.Role {
	if v, ok := claims["role"].(string); ok {
		return domain.Role(v)
	}
	return ""
}