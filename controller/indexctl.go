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
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	qModels, pagination := service.QnA.GetQuestions(page)
	questions := questionsVos(qModels)
	dataModel.Put("Questions", questions)
	dataModel.Put("Pagination", pagination)

	tModels := service.Tag.GetTopTags(6)
	keywords := ""
	for _, t := range tModels {
		keywords += t.Title + ","
	}
	keywords = keywords[:len(keywords)-1]
	dataModel.Put("MetaKeywords", "程序员,编程,代码,问答," + keywords)
	dataModel.Put("MetaDescription", util.Slogan)

	c.HTML(http.StatusOK, "index.html", dataModel)
}
