package myjson

import (
	"encoding/json"
	"log"
	"net/http"
)

func To(input any) ([]byte, error) {
	return json.Marshal(input)
}

func From(source []byte, dest any) error {
	return json.Unmarshal(source, dest)
}

func JsonError(w http.ResponseWriter, status int, msg string) {
	resp, err := To(map[string]interface{}{
		"status":  status,
		"message": msg,
	})
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(status)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err.Error())
	}
}

func JsonResponce(w http.ResponseWriter, status int, msg any) {
	w.Header().Add("Content-Type", "application/json")
	respJSON, err := To(msg)
	if err != nil {
		log.Println(err.Error())
	}

	w.WriteHeader(status)
	_, err = w.Write(respJSON)
	if err != nil {
		log.Println(err.Error())
	}
}
