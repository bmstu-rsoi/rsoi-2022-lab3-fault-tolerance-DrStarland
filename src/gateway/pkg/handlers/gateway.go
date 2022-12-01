package handlers

import (
	"encoding/json"
	"gateway/pkg/models/flights"
	"gateway/pkg/models/tickets"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"go.uber.org/zap"
)

type GatewayHandler struct {
	TicketServiceAddress string
	FlightServiceAddress string
	BonusServiceAddress  string
	Logger               *zap.SugaredLogger
}

func (h *GatewayHandler) GetAllFlights(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	params := r.URL.Query()

	flightsSlice, err := services.GetAllFlightsInfo(h.FlightServiceAddress)
	if err != nil {
		h.Logger.Errorln("failed to get response from flighst service: " + err.Error())
		myjson.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pageParam := params.Get("page")
	if pageParam == "" {
		log.Println("invalid query parameter: (page)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		h.Logger.Errorln("unable to convert the string into int:  " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		h.Logger.Errorln("invalid query parameter (size)")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		h.Logger.Errorln("unable to convert the string into int:  " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	right := page * size
	if len(*flightsSlice) < right {
		right = len(*flightsSlice)
	}

	flightsStripped := (*flightsSlice)[(page-1)*size : right]
	cars := flights.FlightsLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(flightsStripped),
		Items:         &flightsStripped,
	}

	w.Header().Add("Content-Type", "application/json")

	myjson.JsonResponce(w, http.StatusOK, cars)
}

func (h *GatewayHandler) GetUserTickets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketsInfo, err := services.GetUserTickets(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	if err != nil {
		h.Logger.Errorln("Failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticketsInfo); err != nil {
		h.Logger.Errorln("Failed to encode response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) CancelTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		h.Logger.Info("Username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := services.CancelTicket(
		h.TicketServiceAddress,
		h.BonusServiceAddress,
		ps.ByName("ticketUID"),
		username,
	)

	if err != nil {
		h.Logger.Errorln("Failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GatewayHandler) GetUserTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		h.Logger.Errorln("Username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// h.Logger.Infoln("Where is nil 1? ", username)

	ticketUID := ps.ByName("ticketUID")
	// h.Logger.Infoln("Where is nil 2? ", ticketUID)

	ticketsInfo, err := services.GetUserTickets(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	// h.Logger.Infoln("Where is nil 3?", ticketsInfo, err)

	if err != nil {
		h.Logger.Errorln("Failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// h.Logger.Infoln("Where is nil 4?")
	var ticketInfo *tickets.TicketInfo
	for _, ticket := range *ticketsInfo {
		if ticket.TicketUID == ticketUID {
			ticketInfo = &ticket
		}
	}
	// h.Logger.Infoln("Where is nil 5? ", ticketInfo)
	// h.Logger.Info(ticketUID, ticketInfo)
	if ticketInfo == nil {
		myjson.JsonError(w, http.StatusNotFound, "Ticket not found")
		return
	}
	// h.Logger.Infoln("Where is nil 6? ", ticketInfo)
	myjson.JsonResponce(w, http.StatusOK, ticketInfo)
}

func (h *GatewayHandler) BuyTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		h.Logger.Errorln("Username header is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// h.Logger.Infoln("CRINGE1 " + username)
	var ticketInfo tickets.BuyTicketInfo

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.Infoln(err.Error())
	}
	r.Body.Close()
	// h.Logger.Infoln("CRINGE2 " + string(body))

	err = myjson.From(body, &ticketInfo)
	if err != nil {
		h.Logger.Errorln("failed to decode post request: " + err.Error())
		myjson.JsonError(w, http.StatusBadRequest, "failed to decode post request: "+err.Error())
		return
	}
	// h.Logger.Infoln("CRINGE3 ", ticketInfo)

	tickets, err := services.BuyTicket(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		h.BonusServiceAddress,
		username,
		&ticketInfo,
	)

	// h.Logger.Infoln("CRINGE4 ", *tickets)

	if err != nil {
		// h.Logger.Errorln("failed to get response: " + err.Error())
		myjson.JsonError(w, http.StatusBadRequest, "failed to get response: "+err.Error())
		return
	}

	// h.Logger.Debugln("CRINGE4 ", *tickets)

	myjson.JsonResponce(w, http.StatusOK, tickets)
}

func (h *GatewayHandler) GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		myjson.JsonError(w, http.StatusBadRequest, "Username header is empty")
		return
	}

	userInfo, err := services.GetUserInfo(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, "Failed to get response: "+err.Error())
		return
	}

	myjson.JsonResponce(w, http.StatusOK, userInfo)
}

func (h *GatewayHandler) GetPrivilege(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		log.Printf("Username header is empty\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privilegeInfo, err := services.GetUserPrivilege(
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		h.Logger.Errorln("Failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(privilegeInfo); err != nil {
		log.Printf("Failed to encode response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
