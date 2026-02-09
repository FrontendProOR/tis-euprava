package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"tis-euprava/mup-gradjani/internal/service"
)

type updateStatusDTO struct {
	Status string `json:"status"`
}

// /api/requests/{id}/status -> PATCH
func UpdateRequestStatus(svc *service.RequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPatch {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !strings.HasSuffix(r.URL.Path, "/status") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/requests/"), "/status")

		var dto updateStatusDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		processedAt, err := svc.UpdateStatus(id, dto.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":          id,
			"status":      dto.Status,
			"processedAt": processedAt,
		})
	}
}
