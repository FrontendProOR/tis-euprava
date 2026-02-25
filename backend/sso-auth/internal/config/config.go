package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	ServiceName string
	HTTPPort    string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTIssuer       string
	JWTSecret       string
	AccessTTLMin    int
	RefreshTTLDays  int
	CORSAllowOrigin string
}

func LoadConfig() *Config {
	return &Config{
		ServiceName: getEnv("SERVICE_NAME", "service"),
		HTTPPort:    getEnv("HTTP_PORT", "8081"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),

		JWTIssuer:       getEnv("JWT_ISSUER", "tis-euprava-sso"),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		AccessTTLMin:    getEnvInt("ACCESS_TTL_MIN", 15),
		RefreshTTLDays:  getEnvInt("REFRESH_TTL_DAYS", 14),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "http://localhost:5173"),
	}
}

func OpenDB(cfg *Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
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

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}