// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"
	"strconv"

	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/service"
	"github.com/gin-gonic/gin"
)

func showQuestionAction(c *gin.Context) {
	dataModel := getDataModel(c)

	path := c.Param("path")
	qModel := service.QnA.GetQuestionByPath(path)
	if nil == qModel {
		notFound(c)

		return
	}
	question := questionVo(qModel)
	dataModel.Put("Question", question)
	aModels := service.QnA.GetAnswers(qModel.ID)
	dataModel.Put("Answers", answersVos(aModels))
	dataModel.Put("Title", question.Title+" - "+dataModel.GetStr("Title"))
	dataModel.Put("MetaKeywords", qModel.Tags)
	dataModel.Put("MetaDescription", question.Description)

	c.HTML(http.StatusOK, "question.html", dataModel)
}

func showQuestionAnswerAction(c *gin.Context) {
	dataModel := getDataModel(c)

	path := c.Param("path")
	qModel := service.QnA.GetQuestionByPath(path)
	if nil == qModel {
		notFound(c)

		return
	}
	question := questionVo(qModel)
	dataModel.Put("Question", question)
	aModels := service.QnA.GetAnswers(qModel.ID)
	var answerModels []*model.Answer
	aID, err := strconv.Atoi(c.Param("answerID"))
	if nil == err {
		for _, aModel := range aModels {
			if aModel.ID == uint64(aID) {
				answerModels = append(answerModels, aModel)

				break
			}
		}
	} else {
		answerModels = aModels
	}

	dataModel.Put("Answers", answersVos(answerModels))
	dataModel.Put("Title", question.Title+" - "+dataModel.GetStr("Title"))
	dataModel.Put("MetaKeywords", qModel.Tags)
	dataModel.Put("MetaDescription", question.Description)

	c.HTML(http.StatusOK, "question.html", dataModel)
}
