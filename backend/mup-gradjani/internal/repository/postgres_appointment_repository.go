package repository

import (
	"database/sql"

	"tis-euprava/mup-gradjani/internal/domain"
)

type PostgresAppointmentRepository struct{ db *sql.DB }

func NewPostgresAppointmentRepository(db *sql.DB) *PostgresAppointmentRepository {
	return &PostgresAppointmentRepository{db: db}
}

var _ AppointmentRepository = (*PostgresAppointmentRepository)(nil)

func (r *PostgresAppointmentRepository) Create(a *domain.Appointment) error {
	_, err := r.db.Exec(`INSERT INTO appointments (id, citizen_id, date_time, police_station, status) VALUES ($1,$2,$3,$4,$5)`,
		a.ID, a.CitizenID, a.DateTime, a.PoliceStation, a.Status,
	)
	return err
}

func (r *PostgresAppointmentRepository) FindByCitizenID(citizenID string) ([]domain.Appointment, error) {
	rows, err := r.db.Query(`SELECT id, citizen_id, date_time, police_station, status FROM appointments WHERE citizen_id=$1 ORDER BY date_time DESC`, citizenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Appointment
	for rows.Next() {
		var a domain.Appointment
		if err := rows.Scan(&a.ID, &a.CitizenID, &a.DateTime, &a.PoliceStation, &a.Status); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *PostgresAppointmentRepository) Delete(id string) error {
	res, err := r.db.Exec(`DELETE FROM appointments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
