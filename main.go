// 协慌网 - 专注编程问答汉化 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/controller"
	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

// Logger
var logger *gulu.Logger

// The only one init function in RP.
func init() {
	rand.Seed(time.Now().Unix())

	gulu.Log.SetLevel("info")
	logger = gulu.Log.NewLogger(os.Stdout)

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
