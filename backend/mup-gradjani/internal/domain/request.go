package domain

import "time"

type RequestStatus string

const (
	RequestPending   RequestStatus = "PENDING"
	RequestInReview  RequestStatus = "IN_REVIEW"
	RequestApproved  RequestStatus = "APPROVED"
	RequestRejected  RequestStatus = "REJECTED"
	RequestCompleted RequestStatus = "COMPLETED"
)

type ServiceRequest struct {
	ID          string        `json:"id" db:"id"`
	CitizenID   string        `json:"citizen_id" db:"citizen_id"`
	Type        string        `json:"type" db:"type"`
	Status      RequestStatus `json:"status" db:"status"`
	SubmittedAt time.Time     `json:"submitted_at" db:"submitted_at"`
	ProcessedAt *time.Time    `json:"processed_at" db:"processed_at"`
}
