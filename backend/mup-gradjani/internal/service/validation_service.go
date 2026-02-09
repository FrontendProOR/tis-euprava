package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"tis-euprava/mup-gradjani/internal/domain"
)

var (
	ErrValidation = errors.New("validation error")
)

type ValidationService struct{}

func NewValidationService() *ValidationService { return &ValidationService{} }

var jmbgRe = regexp.MustCompile(`^\d{13}$`)

func (v *ValidationService) ValidateCitizen(c *domain.Citizen) error {
	if c == nil {
		return errors.New("citizen is nil")
	}
	c.JMBG = strings.TrimSpace(c.JMBG)
	c.FirstName = strings.TrimSpace(c.FirstName)
	c.LastName = strings.TrimSpace(c.LastName)

	if !jmbgRe.MatchString(c.JMBG) {
		return errors.New("jmbg must be 13 digits")
	}
	if c.FirstName == "" || c.LastName == "" {
		return errors.New("firstName and lastName are required")
	}
	if c.DateOfBirth.IsZero() {
		return errors.New("dateOfBirth is required")
	}
	if c.DateOfBirth.After(time.Now().Add(24 * time.Hour)) {
		return errors.New("dateOfBirth cannot be in the future")
	}
	// Address basic validation (optional but useful)
	if strings.TrimSpace(c.Address.City) == "" {
		return errors.New("address.city is required")
	}
	return nil
}

func (v *ValidationService) ValidateCreateRequest(citizenID, reqType string) error {
	citizenID = strings.TrimSpace(citizenID)
	reqType = strings.TrimSpace(reqType)
	if citizenID == "" {
		return errors.New("citizenId is required")
	}
	if reqType == "" {
		return errors.New("type is required")
	}
	return nil
}

func (v *ValidationService) ValidatePayment(requestID string, amount float64) error {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return errors.New("requestId is required")
	}
	if amount <= 0 {
		return errors.New("amount must be > 0")
	}
	return nil
}
