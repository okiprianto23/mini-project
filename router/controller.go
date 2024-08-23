package router

import (
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"main-xyz/config"
	"net/http"
)

var router *mux.Router

func InitHttpService(controller *HTTPController) {
	router = mux.NewRouter()
	//for api
	controller.Router(router)

	//DEFAULT URL
	controller.HandleFunc(NewHandleFuncParam("/", helloWorld, http.MethodGet))

	//router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Hello World!"))
	//})

	router.Use(MiddlewareCustom)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func StartService(
	config config.Configuration,
	logger *config.LoggerCustom,
) {
	logger.Logger.Info("HTTP Server Start.",
		zap.String("action", "server.start"),
		zap.Int("port", config.Server.Port))
	// Jalankan server HTTPS
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", config.Server.Port), router)
	if err != nil {
		logger.Logger.Fatal("HTTP Server Stopped.",
			zap.String("action", "server.start"),
			zap.Int("port", config.Server.Port))
	}
}
