package domain

import "time"

// Citizen je agregatni root
type Citizen struct {
	ID             string    `json:"id" db:"id"`
	JMBG           string    `json:"jmbg" db:"jmbg"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	DateOfBirth    time.Time `json:"date_of_birth" db:"date_of_birth"`
	Email          string    `json:"email" db:"email"`
	PhoneNumber    string    `json:"phone_number" db:"phone_number"`
	Address        Address   `json:"address" db:"address"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
