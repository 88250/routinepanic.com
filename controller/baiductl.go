// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package controller

import (
	"net/http"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func submitURL(c *gin.Context) {
	for i := 1; i < 256; i++ {
		questions, pagination := service.QnA.GetQuestions(i)
		if i > pagination.LastPageNum {
			logger.Infof("submit completed [p=%d]", pagination.LastPageNum)

			break
		}

		urls := ""
		for _, question := range questions {
			urls += util.Conf.Server + "/questions/" + question.Path + "\n"
		}

		util.PushBaidu(urls)
	}

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}
