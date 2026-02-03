package repository

import (
	"database/sql"
	"encoding/json"

	"tis-euprava/mup-gradjani/internal/domain"
)

type PostgresCitizenRepository struct {
	db *sql.DB
}

func NewPostgresCitizenRepository(db *sql.DB) *PostgresCitizenRepository {
	return &PostgresCitizenRepository{db: db}
}

func (r *PostgresCitizenRepository) Create(c *domain.Citizen) error {
	addressJSON, _ := json.Marshal(c.Address)

	query := `
		INSERT INTO citizens
		(id, jmbg, first_name, last_name, date_of_birth, email, phone_number, address, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err := r.db.Exec(
		query,
		c.ID,
		c.JMBG,
		c.FirstName,
		c.LastName,
		c.DateOfBirth,
		c.Email,
		c.PhoneNumber,
		addressJSON,
		c.CreatedAt,
	)

	return err
}

func (r *PostgresCitizenRepository) FindByID(id string) (*domain.Citizen, error) {
	var c domain.Citizen
	var addressJSON []byte

	query := `
		SELECT id, jmbg, first_name, last_name, date_of_birth,
		       email, phone_number, address, created_at
		FROM citizens WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&c.ID,
		&c.JMBG,
		&c.FirstName,
		&c.LastName,
		&c.DateOfBirth,
		&c.Email,
		&c.PhoneNumber,
		&addressJSON,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(addressJSON, &c.Address)
	return &c, nil
}

func (r *PostgresCitizenRepository) FindByJMBG(jmbg string) (*domain.Citizen, error) {
	var c domain.Citizen
	var addressJSON []byte

	query := `SELECT id, jmbg, first_name, last_name, date_of_birth,
	          email, phone_number, address, created_at
	          FROM citizens WHERE jmbg = $1`

	err := r.db.QueryRow(query, jmbg).Scan(
		&c.ID,
		&c.JMBG,
		&c.FirstName,
		&c.LastName,
		&c.DateOfBirth,
		&c.Email,
		&c.PhoneNumber,
		&addressJSON,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(addressJSON, &c.Address)
	return &c, nil
}
