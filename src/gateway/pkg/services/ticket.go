package services

import (
	"bytes"
	"fmt"
	"gateway/pkg/models/privilege"
	"gateway/pkg/models/tickets"
	"gateway/pkg/myjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func BuyTicket(tAddr, fAddr, bAddr, username string, info *tickets.BuyTicketInfo) (*tickets.PurchaseTicketInfo, error) {
	flight, err := GetFlight(fAddr, info.FlightNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to get flights: %s", err)
	}

	airportFrom, err := GetAirport(fAddr, flight.FromAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	airportTo, err := GetAirport(fAddr, flight.ToAirportId)
	if err != nil {
		return nil, fmt.Errorf("Failed to get airport: %s", err)
	}

	moneyPaid := flight.Price
	bonusesPaid := 0
	diff := int(float32(info.Price) * 0.1)
	optype := "FILL_IN_BALANCE"

	if info.PaidFromBalance {
		if info.Price > 0 {
			bonusesPaid = 0
		} else {
			bonusesPaid = info.Price
		}

		moneyPaid = moneyPaid - bonusesPaid
		diff = -bonusesPaid
		optype = "DEBIT_THE_ACCOUNT"
	}

	uid, err := CreateTicket(tAddr, username, info.FlightNumber, flight.Price)
	if err != nil {
		return nil, fmt.Errorf("Failed to create ticket: %s", err)
	}

	if !info.PaidFromBalance {
		if err := CreatePrivilege(bAddr, username, diff); err != nil {
			return nil, fmt.Errorf("Failed to get privilege info: %s", err)
		}
	}

	err = CreatePrivilegeHistoryRecord(bAddr, uid, flight.Date, optype, 1, diff)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bonus history record: %s", err)
	}

	purchaseInfo := tickets.PurchaseTicketInfo{
		TicketUID:     uid,
		FlightNumber:  info.FlightNumber,
		FromAirport:   fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
		ToAirport:     fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
		Date:          flight.Date,
		Price:         flight.Price,
		PaidByMoney:   moneyPaid,
		PaidByBonuses: bonusesPaid,
		Status:        "PAID",
		Privilege: &privilege.PrivilegeShortInfo{
			Balance: diff,
			Status:  "GOLD",
		},
	}

	return &purchaseInfo, nil
}

func GetTicketsByUsername(ticketsServiceAddress, username string) (*[]tickets.Ticket, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tickets/%s", ticketsServiceAddress, username)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Failed to create an http request")
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to flight service: %w", err)
	}

	tickets := &[]tickets.Ticket{}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	if err = myjson.From(body, tickets); err != nil {
		log.Println(string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tickets, nil
}

func CreateTicket(ticketsServiceAddress, username, flightNumber string, price int) (string, error) {
	requestURL := fmt.Sprintf("%s/api/v1/tickets", ticketsServiceAddress)

	uid := uuid.New().String()
	ticket := &tickets.Ticket{
		TicketUID:    uid,
		FlightNumber: flightNumber,
		Status:       "PAID",
		Username:     username,
		Price:        price,
	}
	data, err := myjson.To(ticket)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		log.Println("Failed to create an http request")
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed request to flights service: %w", err)
	}
	res.Body.Close()

	return uid, nil
}

func CancelTicket(ticketServiceAddress, bonusServiceAddress, ticketUID, username string) error {
	requestURL := fmt.Sprintf("%s/api/v1/tickets/%s", ticketServiceAddress, ticketUID)
	req, err := http.NewRequest(http.MethodDelete, requestURL, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Name", username)
	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to tickets service: %w", err)
	}
	res.Body.Close()
	log.Println("Delete ticket ", res.StatusCode)
	return nil
}
