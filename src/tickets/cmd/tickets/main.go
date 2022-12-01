package main

import (
	"log"
	"net/http"
	"os"
	"tickets/pkg/database"
	"tickets/pkg/handlers"
	"tickets/pkg/models/ticket"

	"github.com/julienschmidt/httprouter"

	mid "tickets/pkg/middleware"

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

	repoTicket := ticket.NewPostgresRepo(db)

	ticketHandler := &handlers.TicketsHandler{
		Logger:      logger,
		TicketsRepo: repoTicket,
	}

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok == true {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	router.POST("/api/v1/tickets", mid.AccessLog(ticketHandler.BuyTicket, logger))
	router.GET("/api/v1/tickets/:username", mid.AccessLog(ticketHandler.GetTicketsByUsername, logger))
	router.DELETE("/api/v1/tickets/:ticketUID", mid.AccessLog(ticketHandler.DeleteTicket, logger))

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
