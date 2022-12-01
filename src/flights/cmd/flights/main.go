package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"flights/pkg/database"

	"flights/pkg/handlers"

	mid "flights/pkg/middleware"
	"flights/pkg/models/airport"
	"flights/pkg/models/flight"

	"go.uber.org/zap"
)

func HealthOK(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	db, err := database.CreateConnection()
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	repoFlight := flight.NewPostgresRepo(db)
	repoAirport := airport.NewPostgresRepo(db)

	allHandler := &handlers.FlightsHandler{
		Logger:      logger,
		FlightRepo:  repoFlight,
		AirportRepo: repoAirport,
	}

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok == true {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	router.GET("/api/v1/flights", mid.AccessLog(allHandler.GetAllFlights, logger))
	router.GET("/api/v1/flights/:flightNumber", mid.AccessLog(allHandler.GetFlight, logger))
	router.GET("/api/v1/airport/:airportID", mid.AccessLog(allHandler.GetAirport, logger))
	router.GET("/manage/health", HealthOK)

	ServerAddress := os.Getenv("PORT")
	if ServerAddress == "" || ServerAddress == ":80" {
		ServerAddress = ":8080"
	}

	logger.Infow("starting server",
		"type", "START",
		"addr", ServerAddress,
	)
	err = http.ListenAndServe(ServerAddress, router)
	if err != nil {
		log.Panicln(err.Error())
	}
}
