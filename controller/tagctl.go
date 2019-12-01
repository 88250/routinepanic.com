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

	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func showTagAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	tagTitle := c.Param("tag")
	tagTitle = tagTitle[1:]
	tagModel := service.Tag.GetTagByTitle(tagTitle)
	if nil == tagModel {
		notFound(c)

		return
	}

	qModels, pagination := service.QnA.GetTagQuestions(tagModel.ID, page)
	questions := questionsVos(qModels)
	dataModel.Put("Questions", questions)
	dataModel.Put("Pagination", pagination)

	dataModel.Put("MetaKeywords", util.MetaKeywords+","+tagModel.Title)

	c.HTML(http.StatusOK, "tag.html", dataModel)
}
