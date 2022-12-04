package handlers

import (
	"gateway/pkg/models/flights"
	"gateway/pkg/models/privilege"
	"gateway/pkg/models/tickets"
	"gateway/pkg/myjson"
	"gateway/pkg/services"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"

	"go.uber.org/zap"
)

type GatewayHandler struct {
	TicketServiceAddress string
	FlightServiceAddress string
	BonusServiceAddress  string
	Logger               *zap.SugaredLogger
}

func (h *GatewayHandler) checkUserHeader(r *http.Request) (string, bool) {
	username := r.Header.Get("X-User-Name")
	if username == "" {
		h.Logger.Errorln("Username header is empty")
		return username, false
	}
	return username, true
}

func (h *GatewayHandler) GetAllFlights(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	params := r.URL.Query()

	flightsSlice, err := services.GetAllFlightsInfo(h.FlightServiceAddress)
	if err != nil {
		h.Logger.Errorln("failed to get response from flight service: " + err.Error())
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
	result := flights.FlightsLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(flightsStripped),
		Items:         &flightsStripped,
	}

	myjson.JsonResponce(w, http.StatusOK, result)
}

func (h *GatewayHandler) GetUserTickets(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketsInfo, err := services.GetUserTickets(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	if err != nil {
		h.Logger.Errorln("failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	myjson.JsonResponce(w, http.StatusOK, ticketsInfo)

	w.WriteHeader(http.StatusOK)
}

func (h *GatewayHandler) CancelTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ticketUID := ps.ByName("ticketUID")

	ticketsInfo, err := services.GetUserTickets(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		username,
	)

	// h.Logger.Infoln("Where is nil 3?", ticketsInfo, err)

	if err != nil {
		h.Logger.Errorln("failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// h.Logger.Infoln("Where is nil 4?")
	var ticketInfo *tickets.TicketInfo
	for _, ticket := range *ticketsInfo {
		if ticket.TicketUID == ticketUID {
			ticketInfo = &ticket
			break
		}
	}
	// h.Logger.Infoln("Where is nil 5? ", ticketInfo)
	// h.Logger.Info(ticketUID, ticketInfo)
	if ticketInfo == nil {
		myjson.JsonError(w, http.StatusNotFound, "ticket not found")
		return
	}

	err = services.CancelTicket(
		h.TicketServiceAddress,
		h.BonusServiceAddress,
		ticketUID,
		username,
	)

	if err != nil {
		h.Logger.Errorln("failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userPrivelege, err := services.GetUserPrivilege(h.BonusServiceAddress, username)
	if err != nil {
		h.Logger.Errorln(err.Error())
		w.WriteHeader(http.StatusNoContent)

		go _cancelTail(h.BonusServiceAddress, ticketUID, username, ticketInfo.Price)
		return
	}

	var bonusRecord *privilege.PrivilegeHistory
	for _, record := range *userPrivelege.History {
		if record.TicketUID == ticketUID {
			bonusRecord = &record
		}
	}

	if bonusRecord == nil {
		h.Logger.Errorln("ООООЙ")
	} else {
		h.Logger.Infoln(*bonusRecord)
	}

	if bonusRecord.OperationType == "DEBIT_THE_ACCOUNT" {
		newBalance := userPrivelege.Balance - (ticketInfo.Price / 10)
		h.Logger.Infoln(userPrivelege.Balance, newBalance)
		err = services.UpdatePrivilege(h.BonusServiceAddress, username, newBalance)
	} else if bonusRecord.OperationType == "FILL_IN_BALANCE" {
		newBalance := userPrivelege.Balance - bonusRecord.BalanceDiff
		h.Logger.Infoln(userPrivelege.Balance, newBalance)
		err = services.UpdatePrivilege(h.BonusServiceAddress, username, newBalance)
	}

	if err != nil {
		h.Logger.Errorln("Cancel ticket: ", err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}

func _cancelTail(address, ticketUID, username string, price int) {
	time.Sleep(5 * time.Second)
	userPrivelege, err := services.GetUserPrivilege(address, username)
	if err != nil {
		go _cancelTail(address, ticketUID, username, price)
		return
	}

	var bonusRecord *privilege.PrivilegeHistory
	for _, record := range *userPrivelege.History {
		if record.TicketUID == ticketUID {
			bonusRecord = &record
		}
	}

	if bonusRecord != nil {
		if bonusRecord.OperationType == "DEBIT_THE_ACCOUNT" {
			newBalance := userPrivelege.Balance - (price / 10)
			err = services.UpdatePrivilege(address, username, newBalance)
		} else if bonusRecord.OperationType == "FILL_IN_BALANCE" {
			newBalance := userPrivelege.Balance - bonusRecord.BalanceDiff
			err = services.UpdatePrivilege(address, username, newBalance)
		}

		if err != nil {
			go _cancelTail(address, ticketUID, username, price)
		}
	}
}

func (h *GatewayHandler) GetUserTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
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
		h.Logger.Errorln("failed to get response: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// h.Logger.Infoln("Where is nil 4?")
	var ticketInfo *tickets.TicketInfo
	for _, ticket := range *ticketsInfo {
		if ticket.TicketUID == ticketUID {
			ticketInfo = &ticket
			break
		}
	}
	// h.Logger.Infoln("Where is nil 5? ", ticketInfo)
	// h.Logger.Info(ticketUID, ticketInfo)
	if ticketInfo == nil {
		myjson.JsonError(w, http.StatusNotFound, "ticket not found")
		return
	}
	// h.Logger.Infoln("Where is nil 6? ", ticketInfo)
	myjson.JsonResponce(w, http.StatusOK, ticketInfo)
}

func (h *GatewayHandler) BuyTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
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
		h.Logger.Errorln("failed to get response: " + err.Error())
		myjson.JsonError(w, http.StatusServiceUnavailable, "Bonus Service unavailable")
		return
	}
	// h.Logger.Debugln("CRINGE4 ", *tickets)
	myjson.JsonResponce(w, http.StatusOK, tickets)
}

func (h *GatewayHandler) GetUserInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, err := services.GetUserInfo(
		h.TicketServiceAddress,
		h.FlightServiceAddress,
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		if err != http.ErrServerClosed {
			h.Logger.Errorln("failed to get response: " + err.Error())
			myjson.JsonError(w, http.StatusInternalServerError, "failed to get response: "+err.Error())
			return
		}

		pseudoUserInfo := struct {
			Privilege   string                `json:"privilege"`
			TicketsInfo *[]tickets.TicketInfo `json:"tickets"`
		}{
			TicketsInfo: userInfo.TicketsInfo,
		}

		myjson.JsonResponce(w, http.StatusOK, pseudoUserInfo)
		return
	}

	myjson.JsonResponce(w, http.StatusOK, userInfo)
}

func (h *GatewayHandler) GetPrivilege(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, isUsername := h.checkUserHeader(r)
	if !isUsername {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	privilegeInfo, err := services.GetUserPrivilege(
		h.BonusServiceAddress,
		username,
	)

	if err != nil {
		h.Logger.Errorln("failed to get response: " + err.Error())
		myjson.JsonError(w, http.StatusServiceUnavailable, "Bonus Service unavailable")
		return
	}

	myjson.JsonResponce(w, http.StatusOK, privilegeInfo)
}
