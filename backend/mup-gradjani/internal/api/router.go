package api

import (
	"net/http"

	"tis-euprava/mup-gradjani/internal/api/handlers"
	"tis-euprava/mup-gradjani/internal/service"
)

func RegisterRoutes(
	mux *http.ServeMux,
	reqSvc *service.RequestService,
	citizenSvc *service.CitizenService,
	apptSvc *service.AppointmentService,
	paySvc *service.PaymentService,
	certSvc *service.CertificateService,
) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.Handle("/api/requests", handlers.Requests(reqSvc, paySvc))
	mux.Handle("/api/requests/", RequestRouter(reqSvc, paySvc, certSvc))

	mux.Handle("/api/citizens", handlers.Citizens(citizenSvc))
	mux.Handle("/api/citizens/", handlers.CitizenByID(citizenSvc))

	mux.Handle("/api/appointments", handlers.Appointments(apptSvc))
	mux.Handle("/api/appointments/", handlers.AppointmentByID(apptSvc))

	mux.Handle("/api/payments", handlers.Payments(paySvc))

}
