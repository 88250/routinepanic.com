package controller

import (
	"strconv"

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/spider"
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

	if importing {
		return
	}

	importing = true
	qnas := spider.StackOverflow.ParseQuestionsByVotes(page, 50)

	for _, qna := range qnas {
		qna.Question.TitleZhCN = service.Translation.Translate(qna.Question.TitleEnUS, "text")
		qna.Question.ContentZhCN = service.Translation.Translate(qna.Question.ContentEnUS, "html")
		for _, a := range qna.Answers {
			a.ContentZhCN = service.Translation.Translate(a.ContentEnUS, "html")
		}
	}

	if err := service.QnA.AddAll(qnas); nil != err {
		logger.Errorf("add QnAs failed: " + err.Error())
	}

	importing = false
}
