package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"bonus/pkg/models/privilege"
	"bonus/pkg/myjson"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type BonusHandler struct {
	Logger    *zap.SugaredLogger
	BonusRepo privilege.Repository
}

func (h *BonusHandler) CreatePrivilegeHistory(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	record := &privilege.PrivilegeHistory{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
	}
	r.Body.Close()

	if err = myjson.From(body, record); err != nil {
		myjson.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.BonusRepo.CreateHistoryRecord(record); err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *BonusHandler) CreatePrivilege(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	record := &privilege.Privilege{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorln(err.Error())
	}
	r.Body.Close()

	if err = myjson.From(body, record); err != nil {
		myjson.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.Logger.Infoln("Trying to create privilege")
	if err = h.BonusRepo.CreatePrivilege(record); err != nil {
		h.Logger.Infoln("chto ne tak " + err.Error())
		oldRecord, _ := h.BonusRepo.GetPrivilegeByUsername(record.Username)
		record.Balance += oldRecord.Balance
		h.Logger.Infoln(oldRecord, record)

		if err = h.BonusRepo.UpdatePrivilege(record); err != nil {
			h.Logger.Infoln("Chto ne tak " + err.Error())
			myjson.JsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BonusHandler) UpdatePrivilege(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	defer func() {
		if err := recover().(error); err != nil {
			h.Logger.Errorln("Recovered in f: " + err.Error())
		}
	}()
	record := &privilege.Privilege{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.Errorln(err.Error())
	}
	r.Body.Close()

	if err = myjson.From(body, record); err != nil {
		myjson.JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.BonusRepo.UpdatePrivilege(record); err != nil {
		h.Logger.Infoln("Chto ne tak " + err.Error())
		myjson.JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BonusHandler) GetPrivilegeByUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	privilege, err := h.BonusRepo.GetPrivilegeByUsername(ps.ByName("username"))
	if err != nil {
		myjson.JsonError(w, http.StatusNotFound, err.Error())
		return
	}

	myjson.JsonResponce(w, http.StatusOK, privilege)
}

func (h *BonusHandler) GetHistoryByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	history, err := h.BonusRepo.GetHistoryById(ps.ByName("privilegeID"))
	if err != nil {
		myjson.JsonError(w, http.StatusInternalServerError, err.Error())
	}

	myjson.JsonResponce(w, http.StatusOK, history)
}
