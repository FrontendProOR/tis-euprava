package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func Requests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			var dto createRequestDTO
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			if dto.CitizenID == "" || dto.Type == "" {
				http.Error(w, "citizenId and type are required", http.StatusBadRequest)
				return
			}

			id := uuid.NewString()
			now := time.Now().UTC()

			_, err := db.Exec(`
				INSERT INTO service_requests (id, citizen_id, type, status, submitted_at)
				VALUES ($1,$2,$3,$4,$5)
			`, id, dto.CitizenID, dto.Type, "SUBMITTED", now)
			if err != nil {
				http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(requestResponse{
				ID:          id,
				CitizenID:   dto.CitizenID,
				Type:        dto.Type,
				Status:      "SUBMITTED",
				SubmittedAt: now,
			})
			return

		case http.MethodGet:
			rows, err := db.Query(`
				SELECT id, citizen_id, type, status, submitted_at, processed_at
				FROM service_requests
				ORDER BY submitted_at DESC
			`)
			if err != nil {
				http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			list := make([]requestResponse, 0)
			for rows.Next() {
				var x requestResponse
				if err := rows.Scan(&x.ID, &x.CitizenID, &x.Type, &x.Status, &x.SubmittedAt, &x.ProcessedAt); err != nil {
					http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
					return
				}
				list = append(list, x)
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
