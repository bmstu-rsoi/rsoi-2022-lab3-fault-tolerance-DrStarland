package ticket

import (
	"database/sql"
	"fmt"
)

type TicketPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *TicketPostgresRepository {
	return &TicketPostgresRepository{DB: db}
}

func (repo *TicketPostgresRepository) GetByUsername(username string) ([]*Ticket, error) {
	// _tickets := make([]Ticket, 0)
	// _rows, _ := repo.DB.Query("SELECT * FROM ticket")
	// defer _rows.Close()
	// for _rows.Next() {
	// 	ticket := Ticket{}
	// 	err := _rows.Scan(
	// 		&ticket.ID,
	// 		&ticket.TicketUID,
	// 		&ticket.Username,
	// 		&ticket.FlightNumber,
	// 		&ticket.Price,
	// 		&ticket.Status)

	// 	if err != nil {
	// 		break
	// 	}

	// 	_tickets = append(_tickets, ticket)
	// }
	// log.Println("DEATH ", _tickets)

	tickets := make([]*Ticket, 0)
	rows, err := repo.DB.Query("SELECT id, ticket_uid, username, flight_number, price, status FROM ticket WHERE username = $1;", username)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		ticket := new(Ticket)
		err = rows.Scan(
			&ticket.ID,
			&ticket.TicketUID,
			&ticket.Username,
			&ticket.FlightNumber,
			&ticket.Price,
			&ticket.Status)

		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %s", err)
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (repo *TicketPostgresRepository) Add(ticket *Ticket) error {
	// _tickets := make([]Ticket, 0)
	// _rows, _ := repo.DB.Query("SELECT * FROM ticket")
	// defer _rows.Close()
	// for _rows.Next() {
	// 	ticket := Ticket{}
	// 	err := _rows.Scan(
	// 		&ticket.ID,
	// 		&ticket.TicketUID,
	// 		&ticket.Username,
	// 		&ticket.FlightNumber,
	// 		&ticket.Price,
	// 		&ticket.Status)

	// 	if err != nil {
	// 		break
	// 	}

	// 	_tickets = append(_tickets, ticket)
	// }
	// log.Println("DEATH ", _tickets)

	_, err := repo.DB.Query(
		"INSERT INTO ticket (ticket_uid, username, flight_number, price, status) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		ticket.TicketUID,
		ticket.Username,
		ticket.FlightNumber,
		ticket.Price,
		ticket.Status,
	)

	return err
}

func (repo *TicketPostgresRepository) Delete(ticketUID string) error {
	_, err := repo.DB.Exec(
		// "DELETE FROM ticket WHERE ticket_uid = $1;",
		"UPDATE ticket SET status = $1 WHERE ticket_uid = $2;",
		"CANCELED",
		ticketUID,
	)
	if err != nil {
		return err
	}

	return err
}
