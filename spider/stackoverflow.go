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

package spider

import (
	"html"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/util"
	"github.com/parnurzeal/gorequest"
)

var StackOverflow = &stackOverflow{}

type stackOverflow struct{}

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

const stackExchangeAPI = "https://api.stackexchange.com"

// QnA represents a question and its answers.
type QnA struct {
	Question *model.Question
	Answers  []*model.Answer
}

func (s *stackOverflow) ParseQuestionsByVotes(page, pageSize int) (ret []*QnA) {
	logger.Info("questions requesting [page=" + strconv.Itoa(page) + ", pageSize=" + strconv.Itoa(pageSize) + "]")
	request := gorequest.New()
	var url = stackExchangeAPI + "/2.2/questions?page=" + strconv.Itoa(page) + "&pagesize=" + strconv.Itoa(pageSize) + "&order=desc&sort=votes&site=stackoverflow&filter=!9Z(-wwYGT"
	data := map[string]interface{}{}
	response, body, errs := request.Set("User-Agent", util.UserAgent).Get(url).Timeout(30 * time.Second).Retry(3, 5*time.Second).EndStruct(&data)
	logger.Info("questions requested [page=" + strconv.Itoa(page) + ", pageSize=" + strconv.Itoa(pageSize) + "]")
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d], response body is [%s]", url, response.StatusCode, body)

		return nil
	}

	qs := data["items"].([]interface{})
	for _, qi := range qs {
		q := qi.(map[string]interface{})
		question := &model.Question{}
		title := q["title"].(string)
		title = html.UnescapeString(title)
		question.TitleEnUS = title
		tis := q["tags"].([]interface{})
		var tags []string
		for _, ti := range tis {
			tags = append(tags, ti.(string))
		}
		question.Tags = strings.Join(tags, ",")
		question.Votes = int(q["score"].(float64))
		question.Views = int(q["view_count"].(float64))
		content := q["body"].(string)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
		doc.Find("pre,code").Each(func(i int, s *goquery.Selection) {
			s.SetAttr("translate", "no")
		})
		question.ContentEnUS, _ = doc.Find("body").Html()
		link := q["link"].(string)
		qId := strconv.Itoa(int(q["question_id"].(float64)))
		question.Path = strings.Split(link, qId+"/")[1]
		question.Source = model.SourceStackOverflow
		question.SourceID = qId
		question.SourceURL = link
		owner := q["owner"].(map[string]interface{})
		authorName := ""
		if nil == owner["display_name"] {
			authorName = "someone"
		} else {
			authorName = owner["display_name"].(string)
		}
		question.AuthorName = authorName
		l := owner["link"]
		if nil != l {
			question.AuthorURL = l.(string)
		}

		answers := s.ParseAnswers(qId)
		qna := &QnA{Question: question, Answers: answers}
		ret = append(ret, qna)

		logger.Infof("parsed voted question [id=%s]", qna.Question.SourceID)
	}

	logger.Infof("parsed voted questions [page=%d]", page)

	return
}

func (s *stackOverflow) ParseAnswers(questionId string) (ret []*model.Answer) {
	logger.Info("answer requesting for question [id=" + questionId + "]")
	request := gorequest.New()
	var url = stackExchangeAPI + "/2.2/questions/" + questionId + "/answers?pagesize=3&order=desc&sort=votes&site=stackoverflow&filter=!9Z(-wzu0T"
	data := map[string]interface{}{}
	response, _, errs := request.Set("User-Agent", util.UserAgent).Get(url).Timeout(30 * time.Second).Retry(3, 5*time.Second).EndStruct(&data)
	logger.Info("answer requested [questionId=" + questionId + "]")
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d]", url, response.StatusCode)

		return nil
	}

	as := data["items"].([]interface{})
	for _, ai := range as {
		a := ai.(map[string]interface{})
		answer := &model.Answer{}
		answer.Votes = int(a["score"].(float64))
		content := a["body"].(string)
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
		doc.Find("pre,code").Each(func(i int, s *goquery.Selection) {
			s.SetAttr("translate", "no")
		})
		answer.ContentEnUS, _ = doc.Find("body").Html()
		answer.Source = model.SourceStackOverflow
		answer.SourceID = strconv.Itoa(int(a["answer_id"].(float64)))
		owner := a["owner"].(map[string]interface{})
		if nil != owner {
			n := owner["display_name"]
			if nil != n {
				answer.AuthorName = n.(string)
			}
			l := owner["link"]
			if nil != l {
				answer.AuthorURL = l.(string)
			}
		}

		ret = append(ret, answer)
	}

	logger.Info("parsed answers for question [id=" + questionId + "]")

	return
}
