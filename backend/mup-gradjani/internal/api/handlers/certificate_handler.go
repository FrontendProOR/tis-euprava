package handlers

import (
	"net/http"
	"strings"

	"tis-euprava/mup-gradjani/internal/service"
)

func CertificateByRequestID(svc *service.CertificateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /api/requests/{id}/certificate
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 3 {
			http.Error(w, "bad path", http.StatusBadRequest)
			return
		}

		requestID := parts[2]
		pdfBytes, err := svc.GenerateCertificate(requestID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdfBytes)
	}
}
