package domain

import "time"

type AppointmentStatus string

const (
	AppointmentScheduled AppointmentStatus = "SCHEDULED"
	AppointmentCancelled AppointmentStatus = "CANCELLED"
)

type Appointment struct {
	ID            string            `json:"id"`
	CitizenID     string            `json:"citizenId"`
	DateTime      time.Time         `json:"dateTime"`
	PoliceStation string            `json:"policeStation"`
	Status        AppointmentStatus `json:"status"`
}
