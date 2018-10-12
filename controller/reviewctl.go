// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"github.com/b3log/routinepanic.com/model"
	"net/http"
	"strconv"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func showReviewAction(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 0, 64)
	rModel := service.Review.GetReviewByID(id)
	if nil == rModel {
		c.Status(http.StatusNotFound)

		return
	}

	dataModel := getDataModel(c)
	review := reviewVo(rModel)
	dataModel.Put("Review", review)

	c.HTML(http.StatusOK, "review.html", dataModel)
}

func showWaitingReviewAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	rModels, pagination := service.Review.GetReviews(model.ReviewStatusWaiting, page)
	reviews := reviewsVos(rModels)
	dataModel.Put("Reviews", reviews)
	dataModel.Put("Pagination", pagination)

	c.HTML(http.StatusOK, "review-waiting.html", dataModel)
}
