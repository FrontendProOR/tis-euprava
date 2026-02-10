package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"

	"tis-euprava/mup-gradjani/internal/api"
	"tis-euprava/mup-gradjani/internal/config"
	"tis-euprava/mup-gradjani/internal/repository"
	"tis-euprava/mup-gradjani/internal/service"
)

func main() {
	log.Println("MUP service started")

	cfg := config.LoadConfig()

	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("Greška pri otvaranju DB konekcije: %v", err)
	}
	defer db.Close()

	log.Println("Povezano na PostgreSQL bazu uspešno!")
	log.Printf("Servis '%s' sluša na portu %s\n", cfg.ServiceName, cfg.HTTPPort)

	// Repositories
	requestRepo := repository.NewPostgresRequestRepository(db)
	citizenRepo := repository.NewPostgresCitizenRepository(db)
	appointmentRepo := repository.NewPostgresAppointmentRepository(db)
	paymentRepo := repository.NewPostgresPaymentRepository(db)
	certificateRepo := repository.NewPostgresCertificateRepository(db)

	// Services
	requestService := service.NewRequestService(requestRepo)
	citizenService := service.NewCitizenService(citizenRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	paymentService := service.NewPaymentService(paymentRepo, requestRepo)
	certificateService := service.NewCertificateService(
		certificateRepo,
		requestRepo,
		paymentRepo,
	)

	// Router
	mux := http.NewServeMux()
	api.RegisterRoutes(
		mux,
		requestService,
		citizenService,
		appointmentService,
		paymentService,
		certificateService,
	)

	// ✅ CORS WRAPPER (OVO TI JE FALILO)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173", // React (Vite)
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	})

	handler := c.Handler(mux)

	// HTTP server
	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           handler, // ⬅️ BITNO
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
