// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"log"
	"os"
	"testing"

	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	util.Conf = &util.Configuration{
		MySQL: "root:@(localhost:3306)/rp?charset=utf8mb4&parseTime=True&loc=Local",
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "test_rp_" + defaultTableName
	}

	ConnectDB()

	log.Println("setup tests")
}

func teardown() {
	DisconnectDB()

	log.Println("teardown tests")
}

func TestAddQuestionsByVotes(t *testing.T) {
	for page := 1; page < 2; page++ {
		qnas := spider.StackOverflow.ParseQuestionsByVotes(page, page)

		//for _, qna := range qnas {
		//	qna.Question.Title = Translation.Translate(qna.Question.Title)
		//	qna.Question.ContentEnUS = Translation.Translate(qna.Question.ContentEnUS)
		//
		//	for _, a := range qna.Answers {
		//		a.ContentEnUS = Translation.Translate(a.ContentEnUS)
		//	}
		//}

		err := QnA.Add(qnas)
		if nil != err {
			t.Errorf("add QnAs failed: " + err.Error())
		}
	}
}
