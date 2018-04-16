// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showQuestionAction(c *gin.Context) {
	dataModel := getDataModel(c)
	c.HTML(http.StatusOK, "question.html", dataModel)
}
