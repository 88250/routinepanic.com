// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"github.com/b3log/routinepanic.com/model"
	"net/http"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func showWaitingReviewAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	rModels, pagination := service.Review.GetReviews(model.ReviewStatusWaiting, page)
	reviews := reviewsVos(rModels)
	dataModel.Put("Reviews", reviews)
	dataModel.Put("Pagination", pagination)

	dataModel.Put("MetaKeywords", util.MetaKeywords)
	dataModel.Put("MetaDescription", util.Slogan)

	c.HTML(http.StatusOK, "review-waiting.html", dataModel)
}
