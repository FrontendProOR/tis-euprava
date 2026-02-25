package api

import (
	"net/http"
	"strings"

	"tis-euprava/mup-gradjani/internal/api/handlers"
	"tis-euprava/mup-gradjani/internal/service"
)

func RequestRouter(svc *service.RequestService, paySvc *service.PaymentService, certSvc *service.CertificateService) http.HandlerFunc {
	getOne := handlers.RequestByID(svc, paySvc)
	patchStatus := handlers.UpdateRequestStatus(svc)
	getCert := handlers.CertificateByRequestID(certSvc)

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/status") {
			patchStatus(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/certificate") {
			getCert(w, r)
			return
		}
		getOne(w, r)
	}
}
