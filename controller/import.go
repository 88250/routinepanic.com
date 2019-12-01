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
	"net/http"
	"strconv"

	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/spider"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

var importing = false

func replenishAnswers(c *gin.Context) {
	service.QnA.ReplenishAnswers()

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}

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
