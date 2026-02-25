package service

import (
	"errors"
	"strings"
)

// Fixed prices (RSD) for demo / faculty project.
// NOTE: We support both "EN" types (ID_CARD/PASSPORT/DRIVER_LICENSE)
// and "SR" types used in DB/UI (LICNA_KARTA/PASOS/VOZACKA).
const (
	PriceIDCard        = 3500.0
	PricePassport      = 4200.0
	PriceDriverLicense = 4000.0
)

func RequiredAmount(requestType string) (float64, error) {
	t := strings.TrimSpace(strings.ToUpper(requestType))

	switch t {
	// ID card
	case "ID_CARD", "LICNA_KARTA", "LICNA":
		return PriceIDCard, nil

	// Passport
	case "PASSPORT", "PASOS":
		return PricePassport, nil

	// Driver license (if you use it)
	case "DRIVER_LICENSE", "VOZACKA", "VOZACKA_DOZVOLA":
		return PriceDriverLicense, nil

	default:
		return 0, errors.New("unknown request type")
	}
}
