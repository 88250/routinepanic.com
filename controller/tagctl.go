// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
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
	dataModel["Questions"] = questions
	dataModel["Pagination"] = pagination

	dataModel["MetaKeywords"] = "程序员,编程,代码,问答," + tagModel.Title
	dataModel["MetaDescription"] = util.Slogan

	c.HTML(http.StatusOK, "tag.html", dataModel)
}