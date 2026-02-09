package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/service"
)

type createCitizenDTO struct {
	JMBG        string         `json:"jmbg"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	DateOfBirth string         `json:"dateOfBirth"` // YYYY-MM-DD
	Email       string         `json:"email,omitempty"`
	PhoneNumber string         `json:"phoneNumber,omitempty"`
	Address     domain.Address `json:"address"`
}

// /api/citizens -> GET(list), POST(create)
func Citizens(svc *service.CitizenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var dto createCitizenDTO
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			dob, err := time.Parse("2006-01-02", strings.TrimSpace(dto.DateOfBirth))
			if err != nil {
				http.Error(w, "dateOfBirth must be YYYY-MM-DD", http.StatusBadRequest)
				return
			}

			cit := &domain.Citizen{
				JMBG:        dto.JMBG,
				FirstName:   dto.FirstName,
				LastName:    dto.LastName,
				DateOfBirth: dob,
				Email:       dto.Email,
				PhoneNumber: dto.PhoneNumber,
				Address:     dto.Address,
			}

			created, err := svc.Create(cit)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(created)
			return

		case http.MethodGet:
			list, err := svc.List()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

// /api/citizens/{id} -> GET
func CitizenByID(svc *service.CitizenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/citizens/")
		if id == "" || strings.Contains(id, "/") {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		cit, err := svc.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(cit)
	}
}
