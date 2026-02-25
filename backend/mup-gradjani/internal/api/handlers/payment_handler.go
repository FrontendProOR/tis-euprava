package handlers

import (
	"encoding/json"
	"net/http"

	"tis-euprava/mup-gradjani/internal/service"
)

type payReq struct {
	RequestID      string  `json:"requestId"`
	RequestIDSnake string  `json:"request_id"`
	Amount         float64 `json:"amount"`
	Reference      string  `json:"reference"`
}

func (p payReq) requestID() string {
	if p.RequestID != "" {
		return p.RequestID
	}
	return p.RequestIDSnake
}

func Payments(svc *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req payReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		p, err := svc.Pay(req.requestID(), req.Amount, req.Reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}
