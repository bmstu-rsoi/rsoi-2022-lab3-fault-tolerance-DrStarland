package flight

import (
	"database/sql"
	"errors"
	"fmt"
)

type FlightPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *FlightPostgresRepository {
	return &FlightPostgresRepository{DB: db}
}

func (repo *FlightPostgresRepository) GetAllFlights() ([]*Flight, error) {
	flights := make([]*Flight, 0)
	rows, err := repo.DB.Query("SELECT * FROM flight;")
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	for rows.Next() {
		f := &Flight{}
		if err := rows.Scan(&f.ID, &f.FlightNumber, &f.Date, &f.FromAirportId, &f.ToAirportId, &f.Price); err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}
		flights = append(flights, f)
	}
	defer rows.Close()

	return flights, nil
}

func (repo *FlightPostgresRepository) GetFlightByNumber(flightNumber string) (*Flight, error) {
	flight := &Flight{}

	err := repo.DB.
		QueryRow("SELECT * FROM flight WHERE flight_number = $1;", flightNumber).
		Scan(
			&flight.ID,
			&flight.FlightNumber,
			&flight.Date,
			&flight.FromAirportId,
			&flight.ToAirportId,
			&flight.Price,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return flight, err
		}
	}

	return flight, nil
}
