package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"tis-euprava/sso-auth/internal/domain"
	"tis-euprava/sso-auth/internal/service"
)

var Auth *service.AuthService
var SSO *service.SSOService

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // CITIZEN / MUP_OFFICER / ZZS_ANALYST
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	role := domain.Role(strings.TrimSpace(strings.ToUpper(req.Role)))
	if role == "" {
		role = domain.RoleCitizen
	}

	u, err := Auth.Register(req.Username, req.Password, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]any{
		"id":       u.ID,
		"username": u.Username,
		"role":     u.Role,
		"active":   u.Active,
	})
}

type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	userAgent := r.Header.Get("User-Agent")
	ip := r.RemoteAddr

	pair, u, err := Auth.Login(req.Username, req.Password, userAgent, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]any{
		"token": pair,
		"user": map[string]any{
			"id":       u.ID,
			"username": u.Username,
			"role":     u.Role,
			"active":   u.Active,
		},
	})
}