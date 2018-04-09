package service

import (
	"testing"
	"github.com/b3log/routinepanic.com/spider"
	"os"
	"log"
	"github.com/b3log/routinepanic.com/util"
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
	ConnectDB()

	log.Println("setup tests")
}

func teardown() {
	DisconnectDB()

	log.Println("teardown tests")
}


func TestParseQuestions(t *testing.T) {
	qnas := spider.StackOverflow.ParseQuestions("https://stackoverflow.com/questions?sort=votes")
	err := QnA.Add(qnas)
	if nil != err {
		t.Errorf("add QnAs failed: " + err.Error())
	}
}
