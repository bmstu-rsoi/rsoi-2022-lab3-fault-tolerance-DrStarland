package services

import (
	"fmt"
	"gateway/pkg/models/airport"
	"gateway/pkg/models/flights"
	"gateway/pkg/myjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetFlight(flightServiceAddress, flightNumber string) (*flights.Flight, error) {
	requestURL := fmt.Sprintf("%s/api/v1/flights/%s", flightServiceAddress, flightNumber)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Failed to create an http request")
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to flight service: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	flight := &flights.Flight{}
	if err = myjson.From(body, flight); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return flight, nil
}

func GetAllFlightsInfo(flightServiceAddress string) (*[]flights.FlightInfo, error) {
	requestURL := fmt.Sprintf("%s/api/v1/flights", flightServiceAddress)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("Failed to create an http request")
		return nil, err
	}
	client := &http.Client{Timeout: 1 * time.Minute}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed request to flight service: %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	flightsSlice := new([]flights.Flight)

	if err = myjson.From(body, flightsSlice); err != nil {
		log.Println("BEDA ", flightsSlice, string(body))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	flightsInfo := make([]flights.FlightInfo, 0)
	for _, flight := range *flightsSlice {
		airportFrom, err := GetAirport(flightServiceAddress, flight.FromAirportId)
		if err != nil {
			return nil, fmt.Errorf("failed to get airport: %s", err)
		}

		airportTo, err := GetAirport(flightServiceAddress, flight.ToAirportId)
		if err != nil {
			return nil, fmt.Errorf("failed to get airport: %s", err)
		}

		fInfo := flights.FlightInfo{
			FlightNumber: flight.FlightNumber,
			FromAirport:  fmt.Sprintf("%s %s", airportFrom.City, airportFrom.Name),
			ToAirport:    fmt.Sprintf("%s %s", airportTo.City, airportTo.Name),
			Date:         flight.Date,
			Price:        flight.Price,
		}

		flightsInfo = append(flightsInfo, fInfo)
	}

	return &flightsInfo, nil
}

func GetAirport(flightServiceAddress string, airportID int) (*airport.Airport, error) {
	requestURL := fmt.Sprintf("%s/api/v1/airport/%d", flightServiceAddress, airportID)
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err.Error())
	}
	res.Body.Close()

	airport := &airport.Airport{}

	if err = myjson.From(body, airport); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return airport, nil
}
