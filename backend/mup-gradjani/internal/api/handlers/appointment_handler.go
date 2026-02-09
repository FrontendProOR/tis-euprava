package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"tis-euprava/mup-gradjani/internal/service"
)

type createAppointmentReq struct {
	CitizenID     string `json:"citizenId"`
	DateTime      string `json:"dateTime"` // ISO
	PoliceStation string `json:"policeStation"`
}

func Appointments(svc *service.AppointmentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req createAppointmentReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			dt, err := time.Parse(time.RFC3339, req.DateTime)
			if err != nil {
				http.Error(w, "invalid dateTime", http.StatusBadRequest)
				return
			}
			created, err := svc.Schedule(req.CitizenID, req.PoliceStation, dt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(created)
		case http.MethodGet:
			citizenID := r.URL.Query().Get("citizenId")
			list, err := svc.ListByCitizen(citizenID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(list)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func AppointmentByID(svc *service.AppointmentService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// očekuje: /api/appointments/{id}
		id := strings.TrimPrefix(r.URL.Path, "/api/appointments/")
		id = strings.Trim(id, "/")

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			err := svc.Cancel(id)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					http.Error(w, "appointment not found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
