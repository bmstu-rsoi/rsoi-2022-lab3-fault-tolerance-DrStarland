package tickets

import "gateway/pkg/models/privilege"

type Ticket struct {
	ID           int    `json:"id"`
	TicketUID    string `json:"ticketUid"`
	Username     string `json:"username"`
	FlightNumber string `json:"flightNumber"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}

type TicketInfo struct {
	TicketUID    string `json:"ticketUid"`
	FlightNumber string `json:"flightNumber"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Date         string `json:"date"`
	Price        int    `json:"price"`
	Status       string `json:"status"`
}

type PurchaseTicketInfo struct {
	TicketUID     string                        `json:"ticketUid"`
	FlightNumber  string                        `json:"flightNumber"`
	FromAirport   string                        `json:"fromAirport"`
	ToAirport     string                        `json:"toAirport"`
	Date          string                        `json:"date"`
	Price         int                           `json:"price"`
	PaidByMoney   int                           `json:"paidByMoney"`
	PaidByBonuses int                           `json:"paidByBonuses"`
	Status        string                        `json:"status"`
	Privilege     *privilege.PrivilegeShortInfo `json:"privilege"`
}

type BuyTicketInfo struct {
	FlightNumber    string `json:"flightNumber"`
	Price           int    `json:"price"`
	PaidFromBalance bool   `json:"paidFromBalance"`
}
