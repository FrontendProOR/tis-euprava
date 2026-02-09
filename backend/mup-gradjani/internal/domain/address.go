package domain

type Address struct {
	Street       string `json:"street" db:"street"`
	Number       string `json:"number" db:"number"`
	City         string `json:"city" db:"city"`
	PostalCode   string `json:"postal_code,omitempty" db:"postal_code"`
	Municipality string `json:"municipality,omitempty" db:"municipality"`
}
