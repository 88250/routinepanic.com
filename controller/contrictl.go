// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func showContriAction(c *gin.Context) {
	dataModel := getDataModel(c)

	dataTypeStr := c.Param("dataType")
	dataType := model.DataTypeQuestion
	if "answers" == dataTypeStr {
		dataType = model.DataTypeAnswer
	}

	dataIdStr := c.Param("id")
	dataId, _ := strconv.ParseUint(dataIdStr, 0, 64)
	if model.DataTypeQuestion == dataType {
		question := service.QnA.GetQuestionByID(dataId)
		if nil == question {
			c.Status(http.StatusNotFound)

			return
		}

		dataModel["Question"] = question

		c.HTML(http.StatusOK, "contri-question.html", dataModel)
	} else {
		answer := service.QnA.GetAnswerByID(dataId)
		if nil == answer {
			c.Status(http.StatusNotFound)

			return
		}

		dataModel["Answer"] = answer

		c.HTML(http.StatusOK, "contri-answer.html", dataModel)
	}
}

func contriAction(c *gin.Context) {
	dataTypeStr := c.Param("dataType")
	dataType := model.DataTypeQuestion
	if "answers" == dataTypeStr {
		dataType = model.DataTypeAnswer
	}
	dataIdStr := c.Param("id")
	dataId, _ := strconv.ParseUint(dataIdStr, 0, 64)
	dataContent := c.PostForm("content")

	questionPath := ""

	if model.DataTypeQuestion == dataType {
		question := service.QnA.GetQuestionByID(dataId)
		if nil == question {
			c.Status(http.StatusNotFound)

			return
		}

		question.ContentZhCN = dataContent
		if err := service.QnA.ContriQuestion(question); nil != err {
			logger.Errorf("contribute to question [%d] failed: %s", dataId, err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}

		questionPath = question.Path
	} else {
		answer := service.QnA.GetAnswerByID(dataId)
		if nil == answer {
			c.Status(http.StatusNotFound)

			return
		}

		question := service.QnA.GetQuestionByID(answer.QuestionID)
		if nil == question {
			c.Status(http.StatusNotFound)

			return
		}
		questionPath = question.Path

		answer.ContentZhCN = dataContent
		if err := service.QnA.ContriAnswer(answer); nil != err {
			logger.Errorf("contribute to answer [%d] failed: %s", dataId, err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}
	}

	c.Redirect(http.StatusSeeOther, util.Conf.Server+"/questions/"+questionPath)
}
