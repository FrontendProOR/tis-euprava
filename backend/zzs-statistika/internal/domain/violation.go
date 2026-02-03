package domain

import "time"

type Violation struct {
	ID         string    `json:"id" db:"id"`
	Type       string    `json:"type" db:"type"`
	Location   string    `json:"location" db:"location"`
	OccurredAt time.Time `json:"occurred_at" db:"occurred_at"`
}
