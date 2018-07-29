// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"log"
	"os"
	"testing"

	"github.com/b3log/routinepanic.com/model"
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

func TestTagQuestions(t *testing.T) {
	var questions []*model.Question

	if err := db.Model(&model.Question{}).Find(&questions).Error; nil != err {
		t.Fatalf("query queestion failed: " + err.Error())
	}

	if err := QnA.TagAll(questions); nil != err {
		t.Errorf("tag questions failed: " + err.Error())
	}
}

func TestReQuestionsByVotes(t *testing.T) {
	for page := 1; page < 10; page++ {
		qnas := spider.StackOverflow.ParseQuestionsByVotes(page, 50)

		if err := QnA.UpdateSourceAll(qnas); nil != err {
			t.Errorf("add QnAs failed: " + err.Error())
		}
	}
}

// Generate data
func TestAddQuestionsByVotes(t *testing.T) {
	for page := 1; page < 10; page++ {
		qnas := spider.StackOverflow.ParseQuestionsByVotes(page, 50)

		for _, qna := range qnas {
			qna.Question.TitleZhCN = Translation.Translate(qna.Question.TitleEnUS, "text")
			qna.Question.ContentZhCN = Translation.Translate(qna.Question.ContentEnUS, "html")
			for _, a := range qna.Answers {
				a.ContentZhCN = Translation.Translate(a.ContentEnUS, "html")
			}
		}

		if err := QnA.AddAll(qnas); nil != err {
			t.Errorf("add QnAs failed: " + err.Error())
		}
	}
}
