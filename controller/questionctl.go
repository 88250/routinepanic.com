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

	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
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
