package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/b3log/routinepanic.com/controller"
	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
)

// Logger
var logger *log.Logger

// Entry point.
func main() {
	service.ConnectDB()

	router := controller.MapRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:" + util.Conf.Port,
		Handler: router,
	}

	handleSignal(server)

	logger.Infof("RP (v%s) is running [%s]", util.Version, util.Conf.Server)
	server.ListenAndServe()
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logger.Infof("got signal [%s], exiting RP now", s)
		if err := server.Close(); nil != err {
			logger.Errorf("server close failed: " + err.Error())
		}

		service.DisconnectDB()

		logger.Infof("RP exited")
		os.Exit(0)
	}()
}
