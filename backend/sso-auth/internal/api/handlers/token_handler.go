package handlers

import (
	"encoding/json"
	"net/http"
)

type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	pair, err := Auth.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	writeJSON(w, pair)
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	claims, err := SSO.ValidateBearer(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]any{
		"valid":  true,
		"claims": claims,
	})
}