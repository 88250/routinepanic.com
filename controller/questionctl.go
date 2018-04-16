// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"

	"github.com/b3log/routinepanic.com/service"
	"github.com/gin-gonic/gin"
)

func showQuestionAction(c *gin.Context) {
	dataModel := getDataModel(c)

	path := c.Request.RequestURI
	path = path[len("/questions/"):]
	qModel := service.QnA.GetQuestionByPath(path)
	if nil == qModel {
		notFound(c)

		return
	}
	dataModel["Question"] = questionVo(qModel)
	aModels := service.QnA.GetAnswers(qModel.ID)
	dataModel["Answers"] = answersVos(aModels)

	c.HTML(http.StatusOK, "question.html", dataModel)
}
