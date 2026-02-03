package main

import (
	"log"
	"net/http"
	"tis-euprava/sso-auth/internal/config"
	"fmt"
)

func main() {
	log.Println("SSO service started")

	cfg := config.LoadConfig()
	db, err := config.OpenDB(cfg)
	if err != nil {
		log.Fatalf("Greška pri otvaranju DB konekcije: %v", err)
	}
	defer db.Close()
	fmt.Printf("Servis '%s' sluša na portu %s\n", cfg.ServiceName, cfg.HTTPPort)
	fmt.Println("Povezano na PostgreSQL bazu uspešno!")

	// Ovdje ide inicijalizacija servisa, rute, handleri, itd.

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
