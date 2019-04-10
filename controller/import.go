// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package controller

import (
	"net/http"
	"strconv"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

var importing = false

func importSO(c *gin.Context) {
	pStr := c.Query("p")
	page, err := strconv.Atoi(pStr)
	if nil != err {
		logger.Errorf("parse page failed: " + err.Error())

		return
	}

	logger.Infof("parsing questions [page="+pStr+", importing=%v]", importing)
	if importing {
		return
	}

	importing = true
	defer func() { importing = false }()
	qnas := spider.StackOverflow.ParseQuestionsByVotes(page, 50)
	logger.Infof("parsed questions [page=" + pStr + "]")
	if err := service.QnA.AddAll(qnas); nil != err {
		logger.Errorf("add QnAs failed: " + err.Error())
	}
	logger.Infof("imported QnAs")

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}

func translate(c *gin.Context) {
	questions := service.QnA.GetUntranslatedQuestions()
	qCnt, aCnt := 0, 0
	for _, q := range questions {
		if "" == q.TitleZhCN {
			q.TitleZhCN = service.Translation.Translate(q.TitleEnUS, "text")
		}
		if "" == q.ContentZhCN {
			q.ContentZhCN = service.Translation.Translate(q.ContentEnUS, "html")
		}

		if err := service.QnA.UpdateQuestion(q); nil != err {
			logger.Errorf("update question failed: " + err.Error())
		}

		logger.Infof("translated a question [" + q.Path + "]")
		qCnt++
	}

	answers := service.QnA.GetUntranslatedAnswers()
	for _, a := range answers {
		if "" == a.ContentZhCN {
			a.ContentZhCN = service.Translation.Translate(a.ContentEnUS, "html")
		}

		if err := service.QnA.UpdateAnswer(a); nil != err {
			logger.Errorf("update answer failed: " + err.Error())
		}

		logger.Infof("translated an answer [%d]", a.ID)
		aCnt++
	}

	logger.Infof("translated questions [%d], answers [%d]", qCnt, aCnt)

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}
