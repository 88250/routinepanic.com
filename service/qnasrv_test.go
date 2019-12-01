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

package service

import (
	"log"
	"os"
	"testing"

	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/spider"
	"github.com/88250/routinepanic.com/util"
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
		MySQL: "root:123456@(localhost:3306)/rp?charset=utf8mb4&parseTime=True&loc=Local",
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
	for page := 3; page < 4; page++ {
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
