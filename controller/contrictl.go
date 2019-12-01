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

package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func tuneHTMLAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses request failed"

		return
	}

	html := arg["html"].(string)
	html = util.TuneHTML(html)

	result.Data = html
}

func showContriAction(c *gin.Context) {
	if !util.IsLoggedIn(c) {
		c.Status(http.StatusUnauthorized)

		return
	}

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

		dataModel.Put("Contri", question)
		dataModel.Put("Type", "Question")
	} else {
		answer := service.QnA.GetAnswerByID(dataId)
		if nil == answer {
			c.Status(http.StatusNotFound)

			return
		}

		dataModel.Put("Contri", answer)
		dataModel.Put("Type", "Answer")
	}

	c.HTML(http.StatusOK, "contri.html", dataModel)
}

func contriAction(c *gin.Context) {
	if !util.IsLoggedIn(c) {
		c.Status(http.StatusUnauthorized)

		return
	}

	dataTypeStr := c.Param("dataType")
	dataType := model.DataTypeQuestion
	if "answers" == dataTypeStr {
		dataType = model.DataTypeAnswer
	}
	dataIdStr := c.Param("id")
	dataId, _ := strconv.ParseUint(dataIdStr, 0, 64)
	dataContent := strings.TrimSpace(c.PostForm("content"))

	session := util.GetSession(c)
	author := service.User.Get(session.UID)

	if model.DataTypeQuestion == dataType {
		question := service.QnA.GetQuestionByID(dataId)
		if nil == question {
			c.Status(http.StatusNotFound)

			return
		}

		title := strings.TrimSpace(c.PostForm("title"))
		question.TitleZhCN = title
		question.ContentZhCN = dataContent
		if err := service.QnA.ContriQuestion(author, question); nil != err {
			logger.Errorf("contribute to question [%d] failed: %s", dataId, err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}
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

		answer.ContentZhCN = dataContent
		if err := service.QnA.ContriAnswer(author, answer); nil != err {
			logger.Errorf("contribute to answer [%d] failed: %s", dataId, err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}
	}

	c.Redirect(http.StatusSeeOther, util.Conf.Server+"/reviews/waiting")
}
