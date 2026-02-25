package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"tis-euprava/mup-gradjani/internal/domain"
	"tis-euprava/mup-gradjani/internal/service"
)

type createRequestDTO struct {
	CitizenID      string `json:"citizenId"`
	CitizenIDSnake string `json:"citizen_id"`
	Type           string `json:"type"`
}

func (d createRequestDTO) citizenID() string {
	if strings.TrimSpace(d.CitizenID) != "" {
		return strings.TrimSpace(d.CitizenID)
	}
	return strings.TrimSpace(d.CitizenIDSnake)
}

type paymentView struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	Status    string  `json:"status"`
	PaidAt    *string `json:"paidAt"`
}

type requestView struct {
	ID          string       `json:"id"`
	CitizenID   string       `json:"citizenId"`
	Type        string       `json:"type"`
	Status      string       `json:"status"`
	SubmittedAt string       `json:"submittedAt"`
	ProcessedAt *string      `json:"processedAt"`
	Price       float64      `json:"price"`
	Paid        bool         `json:"paid"`
	Payment     *paymentView `json:"payment,omitempty"`
}

func toRequestView(r *domain.ServiceRequest, p *domain.Payment) requestView {
	price, _ := service.RequiredAmount(r.Type)

	sub := r.SubmittedAt.Format(timeRFC3339Millis)
	var proc *string
	if r.ProcessedAt != nil {
		s := r.ProcessedAt.Format(timeRFC3339Millis)
		proc = &s
	}

	view := requestView{
		ID:          r.ID,
		CitizenID:   r.CitizenID,
		Type:        r.Type,
		Status:      string(r.Status),
		SubmittedAt: sub,
		ProcessedAt: proc,
		Price:       price,
		Paid:        false,
		Payment:     nil,
	}

	if p != nil {
		view.Paid = (p.Status == domain.PaymentPaid)
		var paidAt *string
		if p.PaidAt != nil {
			s := p.PaidAt.Format(timeRFC3339Millis)
			paidAt = &s
		}
		view.Payment = &paymentView{
			ID:        p.ID,
			Amount:    p.Amount,
			Reference: p.Reference,
			Status:    string(p.Status),
			PaidAt:    paidAt,
		}
	}

	return view
}

const timeRFC3339Millis = "2006-01-02T15:04:05Z07:00"

// /api/requests  -> GET, POST
func Requests(reqSvc *service.RequestService, paySvc *service.PaymentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			var dto createRequestDTO
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			req, err := reqSvc.CreateRequest(dto.citizenID(), dto.Type)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(toRequestView(req, nil))
			return

		case http.MethodGet:
			q := r.URL.Query()
			citizenID := strings.TrimSpace(q.Get("citizenId"))
			status := strings.TrimSpace(q.Get("status"))
			typ := strings.TrimSpace(q.Get("type"))

			var list []domain.ServiceRequest
			var err error
			if citizenID != "" {
				list, err = reqSvc.GetAllByCitizen(citizenID)
			} else {
				list, err = reqSvc.GetAll()
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// optional filters
			if status != "" || typ != "" {
				filtered := make([]domain.ServiceRequest, 0, len(list))
				for _, item := range list {
					if status != "" && string(item.Status) != status {
						continue
					}
					if typ != "" && item.Type != typ {
						continue
					}
					filtered = append(filtered, item)
				}
				list = filtered
			}

			out := make([]requestView, 0, len(list))
			for i := range list {
				item := list[i]
				var pay *domain.Payment
				if p, err := paySvc.GetByRequestID(item.ID); err == nil {
					pay = p
				}
				out = append(out, toRequestView(&item, pay))
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(out)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// /api/requests/{id} -> GET
func RequestByID(reqSvc *service.RequestService, paySvc *service.PaymentService) http.HandlerFunc {
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

		req, err := reqSvc.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var pay *domain.Payment
		if p, err := paySvc.GetByRequestID(req.ID); err == nil {
			pay = p
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(toRequestView(req, pay))
	}
}
