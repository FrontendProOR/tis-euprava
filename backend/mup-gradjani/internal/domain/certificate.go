package domain

import "time"

type ElectronicCertificate struct {
	ID        string    `json:"id"`
	RequestID string    `json:"requestId"`
	IssuedAt  time.Time `json:"issuedAt"`
	PDF       []byte    `json:"-"`
}
