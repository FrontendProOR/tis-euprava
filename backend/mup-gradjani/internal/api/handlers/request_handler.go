package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"tis-euprava/mup-gradjani/internal/service"
)

type createRequestDTO struct {
	CitizenID string `json:"citizenId"`
	Type      string `json:"type"`
}

type requestResponse struct {
	ID          string     `json:"id"`
	CitizenID   string     `json:"citizenId"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	SubmittedAt time.Time  `json:"submittedAt"`
	ProcessedAt *time.Time `json:"processedAt,omitempty"`
}

// /api/requests  -> GET, POST
func Requests(svc *service.RequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			var dto createRequestDTO
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			req, err := svc.CreateRequest(dto.CitizenID, dto.Type)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(req)
			return

		case http.MethodGet:
			list, err := svc.GetAll()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(list)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// /api/requests/{id} -> GET
func RequestByID(svc *service.RequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/api/requests/")
		if id == "" || strings.Contains(id, "/") {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		req, err := svc.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(req)
	}
}
