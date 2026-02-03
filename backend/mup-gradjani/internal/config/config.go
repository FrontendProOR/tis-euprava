package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
)

// Config sadrži sve konfiguracione parametre servisa
type Config struct {
	ServiceName string // naziv servisa, npr "mup-gradjani"
	HTTPPort    string // port na kojem servis sluša, npr "8002"

	// PostgreSQL konekcija
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig učitava konfiguraciju iz environment varijabli
func LoadConfig() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "service"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
	}
}

// OpenDB otvara konekciju sa PostgreSQL bazom
func OpenDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("ne mogu da se povežem na DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping neuspešan: %w", err)
	}

	return db, nil
}

// getEnv vraća vrednost env varijable ili fallback ako nije postavljena
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
