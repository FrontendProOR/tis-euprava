package domain

import "time"

type RequestStatus string

const (
	RequestSubmitted RequestStatus = "SUBMITTED"
	RequestInProcess RequestStatus = "IN_PROCESS"
	RequestApproved  RequestStatus = "APPROVED"
	RequestRejected  RequestStatus = "REJECTED"
)

type ServiceRequest struct {
	ID           string               `json:"id" db:"id"`
	CitizenID    string               `json:"citizen_id" db:"citizen_id"`
	Type         string               `json:"type" db:"type"`
	Status       RequestStatus        `json:"status" db:"status"`
	SubmittedAt  time.Time            `json:"submitted_at" db:"submitted_at"`
	ProcessedAt  *time.Time           `json:"processed_at" db:"processed_at"`
	DocumentData *DocumentRequestData `json:"document_data" db:"document_data"`
}
