// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/proxy"
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

	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1081", nil, proxy.Direct)
	if err != nil {
		log.Fatal("can't connect to the proxy: " + err.Error())
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	http.DefaultClient.Transport = httpTransport
}

func teardown() {
	DisconnectDB()

	log.Println("teardown tests")
}

func TestAddQuestionsByVotes(t *testing.T) {
	for page := 1; page < 10; page++ {
		qnas := spider.StackOverflow.ParseQuestionsByVotes(page, page)

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

func TestReAdd(t *testing.T) {
	var questions []*model.Question

	if err := db.Model(&model.Question{}).Where("`content_zh_cn` = ''").Find(&questions).Error; nil != err {
		t.Fatalf("query queestion failed: " + err.Error())
	}

	for _, q := range questions {
		qna := spider.StackOverflow.ParseQuestion(q.SourceURL)

		qna.Question.TitleZhCN = Translation.Translate(qna.Question.TitleEnUS, "html")
		qna.Question.ContentZhCN = Translation.Translate(qna.Question.ContentEnUS, "html")
		for _, a := range qna.Answers {
			a.ContentZhCN = Translation.Translate(a.ContentEnUS, "html")
		}

		if err := QnA.Add(qna); nil != err {
			t.Errorf("add QnA failed: " + err.Error())
		}
	}
}
