// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showQuestionAction(c *gin.Context) {
	logger.Infof("path " + c.Request.RequestURI)
	dataModel := &DataModel{}
	c.HTML(http.StatusOK, "question.html", dataModel)
}
