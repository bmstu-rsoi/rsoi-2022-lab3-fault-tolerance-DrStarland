package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"tickets/pkg/models/ticket"
	"tickets/pkg/myjson"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type TicketsHandler struct {
	Logger      *zap.SugaredLogger
	TicketsRepo ticket.Repository
}

func (h *TicketsHandler) GetTicketsByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	tickets, err := h.TicketsRepo.GetByUsername(username)
	if err != nil {
		log.Printf("Failed to get ticket: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	myjson.JsonResponce(w, http.StatusOK, tickets)
}

func (h *TicketsHandler) BuyTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Header.Get("Content-Type") != "application/json" {
		myjson.JsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	ticket := &ticket.Ticket{}
	err := myjson.From(body, ticket)

	if err != nil {
		h.Logger.Errorln("STRANDING ", string(body))
		myjson.JsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	if err := h.TicketsRepo.Add(ticket); err != nil {
		log.Printf("Failed to create ticket: %s", err)

		myjson.JsonError(w, http.StatusInternalServerError, "Failed to create ticket: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TicketsHandler) DeleteTicket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ticketUID := ps.ByName("ticketUID")

	if err := h.TicketsRepo.Delete(ticketUID); err != nil {
		h.Logger.Errorln("Failed to create ticket: " + err.Error())

		myjson.JsonError(w, http.StatusInternalServerError, "failed to create ticket: "+err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
