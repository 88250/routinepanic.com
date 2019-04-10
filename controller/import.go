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

	var importQnas []*spider.QnA
	for _, qna := range qnas {
		if service.QnA.Translated(qna) {
			continue
		}

		qna.Question.TitleZhCN = service.Translation.Translate(qna.Question.TitleEnUS, "text")
		qna.Question.ContentZhCN = service.Translation.Translate(qna.Question.ContentEnUS, "html")
		for _, a := range qna.Answers {
			a.ContentZhCN = service.Translation.Translate(a.ContentEnUS, "html")
		}

		logger.Info("translated a QnA [" + qna.Question.TitleEnUS + ", " + qna.Question.TitleZhCN + "]")
		importQnas = append(importQnas, qna)
	}

	if err := service.QnA.AddAll(importQnas); nil != err {
		logger.Errorf("add QnAs failed: " + err.Error())
	}
	logger.Infof("imported QnAs")
	importing = false

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}
