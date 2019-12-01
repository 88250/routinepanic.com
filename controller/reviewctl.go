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
	"fmt"
	"net/http"
	"strconv"

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func ReviewAction(c *gin.Context) {
	if !util.IsLoggedIn(c) {
		c.Status(http.StatusUnauthorized)

		return
	}

	session := util.GetSession(c)
	if util.RoleReviewer != session.URole {
		c.Status(http.StatusUnauthorized)

		return
	}

	result := gulu.Ret.NewResult()
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

	url := ""
	revision := service.QnA.GetRevision(review.RevisionID)
	if model.DataTypeQuestion == revision.DataType {
		question := service.QnA.GetQuestionByID(revision.DataID)
		url = util.Conf.Server + "/questions/" + question.Path
	} else {
		answer := service.QnA.GetAnswerByID(revision.DataID)
		question := service.QnA.GetQuestionByID(answer.QuestionID)
		url = fmt.Sprintf(util.Conf.Server+"/questions/"+question.Path+"/answers/%d", answer.ID)
	}
	result.Data = url

	if passed {
		util.PushBaidu(url)
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
