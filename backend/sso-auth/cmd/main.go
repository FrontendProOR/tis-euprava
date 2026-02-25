package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tis-euprava/sso-auth/internal/api"
	"tis-euprava/sso-auth/internal/api/handlers"
	"tis-euprava/sso-auth/internal/config"
	"tis-euprava/sso-auth/internal/repository"
	"tis-euprava/sso-auth/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	// učitaj .env ako postoji (lokalni dev)
	_ = godotenv.Load()

	log.Println("SSO service starting...")

	cfg := config.LoadConfig()

	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("Greška pri otvaranju DB konekcije: %v", err)
	}
	defer db.Close()

	fmt.Printf("Servis '%s' sluša na portu %s\n", cfg.ServiceName, cfg.HTTPPort)
	fmt.Println("Povezano na PostgreSQL bazu uspešno!")

	// repositories
	userRepo := repository.NewPostgresUserRepository(db)
	sessionRepo := repository.NewPostgresSessionRepository(db)

	// services
	tokenSvc := service.NewTokenService(cfg)
	authSvc := service.NewAuthService(cfg, userRepo, sessionRepo, tokenSvc)
	ssoSvc := service.NewSSOService(tokenSvc)

	// inject into handlers
	handlers.Auth = authSvc
	handlers.SSO = ssoSvc

	// router
	router := api.NewRouter()

	addr := ":" + cfg.HTTPPort
	if os.Getenv("HTTP_PORT") == "" && cfg.HTTPPort == "8081" {
		// ok
	}

	log.Printf("SSO service running at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}