// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"
	"strconv"

	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func ReviewAction(c *gin.Context) {
	if !util.IsLoggedIn(c) {
		c.Status(http.StatusUnauthorized)

		return
	}

	session := util.GetSession(c)
	if "88250" != session.UName || "Vanessa" != session.UName {
		c.Status(http.StatusUnauthorized)

		return
	}

	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	arg := map[string]interface{}{}
	if err := c.BindJSON(&arg); nil != err {
		result.Code = -1
		result.Msg = "parses request failed"

		return
	}

	passed := arg["passed"].(bool)
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 0, 64)

	review := service.Review.GetReviewByID(id)
	review.ReviewerID = session.UID

	if passed {
		if err := service.Review.PassReview(review, arg); nil != err {
			logger.Errorf("pass review failed: %s", err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}
	} else {
		memo := arg["memo"].(string)
		review.Memo = memo
		if err := service.Review.RejectReview(review); nil != err {
			logger.Errorf("reject review failed: %s", err.Error())
			c.Status(http.StatusInternalServerError)

			return
		}
	}
}

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

func showWaitingReviewsAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	rModels, pagination := service.Review.GetReviews(model.ReviewStatusWaiting, page)
	reviews := reviewsVos(rModels)
	dataModel.Put("Reviews", reviews)
	dataModel.Put("Pagination", pagination)
	dataModel.Put("Type", "Waiting")
	c.HTML(http.StatusOK, "reviews.html", dataModel)
}

func showPassedReviewsAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	rModels, pagination := service.Review.GetReviews(model.ReviewStatusPassed, page)
	reviews := reviewsVos(rModels)
	dataModel.Put("Reviews", reviews)
	dataModel.Put("Pagination", pagination)
	dataModel.Put("Type", "Passed")
	c.HTML(http.StatusOK, "reviews.html", dataModel)
}

func showRejectedReviewsAction(c *gin.Context) {
	dataModel := getDataModel(c)
	page := util.GetPage(c)

	rModels, pagination := service.Review.GetReviews(model.ReviewStatusRejected, page)
	reviews := reviewsVos(rModels)
	dataModel.Put("Reviews", reviews)
	dataModel.Put("Pagination", pagination)
	dataModel.Put("Type", "Rejected")
	c.HTML(http.StatusOK, "reviews.html", dataModel)
}
