// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func showIndexAction(c *gin.Context) {
	dataModel := &DataModel{}
	c.HTML(http.StatusOK, "index.html", dataModel)
}
