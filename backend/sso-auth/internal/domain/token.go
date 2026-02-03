package domain

import "time"

type Token struct {
	Value     string    `json:"value" db:"value"`
	UserID    string    `json:"user_id" db:"user_id"`
	IssuedAt  time.Time `json:"issued_at" db:"issued_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}
