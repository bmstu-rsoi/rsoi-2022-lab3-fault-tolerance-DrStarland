package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"gateway/pkg/database"
	"gateway/pkg/handlers"

	mid "gateway/pkg/middleware"

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
	//	repoFlight := flight.NewPostgresRepo(db)
	// allHandler := &handlers.FlightsHandler{
	// 	Logger:     logger,
	// 	FlightRepo: repoFlight,
	// }

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok == true {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	gs := &handlers.GatewayHandler{
		TicketServiceAddress: "http://testum_tickets:8070",
		FlightServiceAddress: "http://testum_flights:8060",
		BonusServiceAddress:  "http://testum_bonus:8050",
		Logger:               logger,
	}

	router.GET("/api/v1/flights", mid.AccessLog(gs.GetAllFlights, logger))
	router.GET("/api/v1/me", mid.AccessLog(gs.GetUserInfo, logger))
	router.GET("/api/v1/tickets", mid.AccessLog(gs.GetUserTickets, logger))
	router.GET("/api/v1/tickets/:ticketUID", mid.AccessLog(gs.GetUserTicket, logger))
	router.POST("/api/v1/tickets", mid.AccessLog(gs.BuyTicket, logger))
	router.DELETE("/api/v1/tickets/:ticketUID", mid.AccessLog(gs.CancelTicket, logger))
	router.GET("/api/v1/privilege", mid.AccessLog(gs.GetPrivilege, logger))

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
