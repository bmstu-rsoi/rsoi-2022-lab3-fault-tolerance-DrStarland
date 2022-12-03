package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"bonus/pkg/database"

	"bonus/pkg/handlers"

	mid "bonus/pkg/middleware"
	"bonus/pkg/models/privilege"

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

	bonusRepo := privilege.NewPostgresRepo(db)

	bonusHandler := &handlers.BonusHandler{
		Logger:    logger,
		BonusRepo: bonusRepo,
	}

	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		log.Println("panicMiddleware is working", r.URL.Path)
		if trueErr, ok := err.(error); ok == true {
			http.Error(w, "Internal server error: "+trueErr.Error(), http.StatusInternalServerError)
		}
	}

	router.POST("/api/v1/bonus", mid.AccessLog(bonusHandler.CreatePrivilegeHistory, logger))
	router.POST("/api/v1/bonus/privilege", mid.AccessLog(bonusHandler.CreatePrivilege, logger))
	router.PUT("/api/v1/bonus/privilege", mid.AccessLog(bonusHandler.UpdatePrivilege, logger))
	router.GET("/api/v1/bonus/:username", mid.AccessLog(bonusHandler.GetPrivilegeByUsername, logger))
	router.GET("/api/v1/bonushistory/:privilegeID", mid.AccessLog(bonusHandler.GetHistoryByID, logger))

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
