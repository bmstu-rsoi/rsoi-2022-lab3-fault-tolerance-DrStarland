package services

import (
	"fmt"
	"gateway/pkg/models/privilege"
	"gateway/pkg/models/tickets"
	"gateway/pkg/models/user"
	"net/http"
)

func GetUserTickets(ticketServiceAddress, flightServiceAddress, username string) (*[]tickets.TicketInfo, error) {
	userTickets, err := GetTicketsByUsername(ticketServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %s", err)
	}

	ticketsInfo := make([]tickets.TicketInfo, 0)
	for _, ticket := range *userTickets {
		flight, err := GetFlight(flightServiceAddress, ticket.FlightNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get flights: %s", err)
		}
		airportFrom, err := GetAirport(flightServiceAddress, flight.FromAirportId)
		if err != nil {
			return nil, fmt.Errorf("failed to get airport: %s", err)
		}
		airportTo, err := GetAirport(flightServiceAddress, flight.ToAirportId)
		if err != nil {
			return nil, fmt.Errorf("failed to get airport: %s", err)
		}

		ticketInfo := tickets.TicketInfo{
			TicketUID:    ticket.TicketUID,
			FlightNumber: ticket.FlightNumber,
			FromAirport:  fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
			ToAirport:    fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
			Date:         flight.Date,
			Price:        ticket.Price,
			Status:       ticket.Status,
		}

		ticketsInfo = append(ticketsInfo, ticketInfo)
	}

	return &ticketsInfo, nil
}

func GetUserInfo(ticketServiceAddress, flightServiceAddress, bonusServiceAddress, username string) (*user.UserInfo, error) {
	ticketsInfo, err := GetUserTickets(ticketServiceAddress, flightServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %s", err)
	}

	userInfo := &user.UserInfo{
		TicketsInfo: ticketsInfo,
		Privilege:   &privilege.PrivilegeShortInfo{},
	}

	userPrivilege, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return userInfo, http.ErrServerClosed
	}

	userInfo.Privilege.Status = userPrivilege.Status
	userInfo.Privilege.Balance = userPrivilege.Balance

	return userInfo, nil
}

func GetUserPrivilege(bonusServiceAddress, username string) (*privilege.PrivilegeInfo, error) {
	privilegeShortInfo, err := GetPrivilegeShortInfo(bonusServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %s", err)
	}

	privilegeHistory, err := GetPrivilegeHistory(bonusServiceAddress, privilegeShortInfo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get privilege info: %s", err)
	}

	privilegeInfo := &privilege.PrivilegeInfo{
		Status:  privilegeShortInfo.Status,
		Balance: privilegeShortInfo.Balance,
		History: privilegeHistory,
	}

	return privilegeInfo, nil
}
