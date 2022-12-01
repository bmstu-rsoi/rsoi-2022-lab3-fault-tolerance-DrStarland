package middleware

import (
	"log"
	"net/http"
)

func PanicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Println("panicMiddleware is working", r.URL.Path)
	if trueErr, ok := err.(error); ok == true {
		http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
	}
}
