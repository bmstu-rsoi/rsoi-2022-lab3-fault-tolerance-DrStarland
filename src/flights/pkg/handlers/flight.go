package handlers

import (
	"net/http"

	"flights/pkg/models/airport"
	"flights/pkg/models/flight"
	"flights/pkg/myjson"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type FlightsHandler struct {
	Logger      *zap.SugaredLogger
	FlightRepo  flight.Repository
	AirportRepo airport.Repository
}

func (h *FlightsHandler) GetAllFlights(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	flights, err := h.FlightRepo.GetAllFlights()
	if err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, "flight service error: "+err.Error())
		return
	}

	myjson.JsonResponce(w, http.StatusOK, flights)
}

func (h *FlightsHandler) GetFlight(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	number := ps.ByName("flightNumber")
	flight, err := h.FlightRepo.GetFlightByNumber(number)
	if err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, "flight service error: "+err.Error())
		return
	}
	myjson.JsonResponce(w, http.StatusOK, flight)
}

func (h *FlightsHandler) GetAirport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("airportID")
	airport, err := h.AirportRepo.GetAirportByID(id)
	if err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, "flight service error: "+err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	myjson.JsonResponce(w, http.StatusOK, airport)
}
