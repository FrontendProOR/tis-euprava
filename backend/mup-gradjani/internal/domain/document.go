package domain

type DocumentType string

const (
	DocumentIDCard   DocumentType = "ID_CARD"
	DocumentPassport DocumentType = "PASSPORT"
	DocumentDriver   DocumentType = "DRIVER_LICENSE"
)

type DocumentRequestData struct {
	Type     DocumentType `json:"type" db:"type"`
	Reason   string       `json:"reason" db:"reason"`
	IsUrgent bool         `json:"is_urgent" db:"is_urgent"`
}
