package domain

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentPaid    PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID        string        `json:"id" db:"id"`
	RequestID string        `json:"request_id" db:"request_id"`
	Amount    float64       `json:"amount" db:"amount"`
	Reference string        `json:"reference" db:"reference"`
	Status    PaymentStatus `json:"status" db:"status"`
	PaidAt    *time.Time    `json:"paid_at" db:"paid_at"`
}
