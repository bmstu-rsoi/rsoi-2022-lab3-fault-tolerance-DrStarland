package airport

import (
	"database/sql"
	"errors"
)

type AirportPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *AirportPostgresRepository {
	return &AirportPostgresRepository{DB: db}
}

// id      SERIAL PRIMARY KEY,
// name    VARCHAR(255),
// city    VARCHAR(255),
// country VARCHAR(255)

func (repo *AirportPostgresRepository) GetAirportByID(airportID string) (*Airport, error) {
	airport := &Airport{}

	err := repo.DB.
		QueryRow("SELECT id, name, city, country FROM airport WHERE id = $1", airportID).
		Scan(
			&airport.ID,
			&airport.Name,
			&airport.City,
			&airport.Country,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return airport, err
		}
	}

	return airport, nil
}
