// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b3log/routinepanic.com/controller"
	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

// Logger
var logger *log.Logger

// The only one init function in RP.
func init() {
	rand.Seed(time.Now().Unix())

	log.SetLevel("info")
	logger = log.NewLogger(os.Stdout)

	util.LoadConf()

	if "dev" == util.Conf.RuntimeMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
}

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
