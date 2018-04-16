// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"
	"strings"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func showIndexAction(c *gin.Context) {
	dataModel := DataModel{}

	page := util.GetPage(c)
	qModels, pagination := service.QnA.GetQuestions(page)
	questions := []*question{}
	for _, qModel := range qModels {
		q := &question{Title: qModel.Title}
		tagStrs := strings.Split(qModel.Tags, ",")
		for _, tagTitle := range tagStrs {
			q.Tags = append(q.Tags, &tag{Title: tagTitle})
		}

		questions = append(questions, q)
	}

	dataModel["Questions"] = questions
	dataModel["Pagination"] = pagination
	c.HTML(http.StatusOK, "index.html", dataModel)
}
