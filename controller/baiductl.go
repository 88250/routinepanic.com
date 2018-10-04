// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"time"

	"github.com/b3log/routinepanic.com/service"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
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
			urls += "https://routinepanic.com/questions/" + question.Path + "\n"
		}

		_, data, errors := gorequest.New().Post("http://data.zz.baidu.com/urls?site=routinepanic.com&token=y82YmzUKXO9JZUAr").
			AppendHeader("User-Agent", "curl/7.12.1").
			AppendHeader("Host", "data.zz.baidu.com").
			AppendHeader("Content-Type", "text/plain").Timeout(10 * time.Second).Send(urls).EndBytes()
		if nil != errors {
			logger.Errorf("submit failed [%s]", errors)

			continue
		}

		logger.Info(string(data))
	}
}
