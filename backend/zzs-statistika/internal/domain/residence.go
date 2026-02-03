package domain

type ResidenceData struct {
	Municipality string `json:"municipality" db:"municipality"`
	Population   int    `json:"population" db:"population"`
}
