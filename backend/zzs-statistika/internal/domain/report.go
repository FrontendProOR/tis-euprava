package domain

import "time"

type Report struct {
	ID        string            `json:"id" db:"id"`
	Title     string            `json:"title" db:"title"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	Data      []AnalyticsResult `json:"data" db:"data"`
}
