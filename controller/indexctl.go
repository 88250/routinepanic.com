// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func showIndexAction(c *gin.Context) {
	dataModel := DataModel{}

	page := util.GetPage(c)
	questions, pagination := service.QnA.GetQuestions(page)
	dataModel["Questions"] = questions
	dataModel["Pagination"] = pagination
	c.HTML(http.StatusOK, "index.html", dataModel)
}
