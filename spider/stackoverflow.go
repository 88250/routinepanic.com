// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package spider

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/util"
	"github.com/parnurzeal/gorequest"
)

var StackOverflow = &stackOverflow{}

type stackOverflow struct{}

// Logger
var logger = log.NewLogger(os.Stdout)

const stackExchangeAPI = "https://api.stackexchange.com"

// QnA represents a question and its answers.
type QnA struct {
	Question *model.Question
	Answers  []*model.Answer
}

func (s *stackOverflow) ParseQuestionsByVotes(page, pageSize int) (ret []*QnA) {
	request := gorequest.New()
	var url = stackExchangeAPI + "/2.2/questions?page=" + strconv.Itoa(page) + "&pagesize=" + strconv.Itoa(pageSize) + "&order=desc&sort=votes&site=stackoverflow&filter=!9Z(-wwYGT"
	data := map[string]interface{}{}
	response, _, errs := request.Set("User-Agent", util.UserAgent).Get(url).Retry(3, 5*time.Second).EndStruct(&data)
	logger.Info("question requested [" + url + "]")
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d]", response.StatusCode)

		return nil
	}

	qs := data["items"].([]interface{})
	for _, qi := range qs {
		q := qi.(map[string]interface{})
		question := &model.Question{}
		question.TitleEnUS = q["title"].(string)
		tis := q["tags"].([]interface{})
		var tags []string
		for _, ti := range tis {
			tags = append(tags, ti.(string))
		}
		question.Tags = strings.Join(tags, ",")
		question.Votes = int(q["score"].(float64))
		question.Views = int(q["view_count"].(float64))
		question.ContentEnUS = q["body"].(string)
		link := q["link"].(string)
		qId := strconv.Itoa(int(q["question_id"].(float64)))
		question.Path = strings.Split(link, qId+"/")[1]
		question.Source = model.SourceStackOverflow
		question.SourceID = qId
		question.SourceURL = link
		owner := q["owner"].(map[string]interface{})
		question.AuthorName = owner["display_name"].(string)
		question.AuthorURL = owner["link"].(string)

		answers := s.ParseAnswers(qId)
		qna := &QnA{Question: question, Answers: answers}
		ret = append(ret, qna)

		logger.Infof("parsed voted question [id=%s]", qna.Question.SourceID)
	}

	logger.Infof("parsed voted questions [page=%d]", page)

	return
}

func (s *stackOverflow) ParseAnswers(questionId string) (ret []*model.Answer) {
	request := gorequest.New()
	var url = stackExchangeAPI + "/2.2/questions/" + questionId + "/answers?pagesize=3&order=desc&sort=votes&site=stackoverflow&filter=!9Z(-wzu0T"
	data := map[string]interface{}{}
	response, _, errs := request.Set("User-Agent", util.UserAgent).Get(url).Retry(3, 5*time.Second).EndStruct(&data)
	logger.Info("answer requested [" + url + "]")
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d]", response.StatusCode)

		return nil
	}

	as := data["items"].([]interface{})
	for _, ai := range as {
		a := ai.(map[string]interface{})
		answer := &model.Answer{}
		answer.Votes = int(a["score"].(float64))
		answer.ContentEnUS = a["body"].(string)
		answer.Source = model.SourceStackOverflow
		answer.SourceID = strconv.Itoa(int(a["answer_id"].(float64)))
		owner := a["owner"].(map[string]interface{})
		answer.AuthorName = owner["display_name"].(string)
		l := owner["link"]
		if nil != l {
			answer.AuthorURL = l.(string)
		}

		ret = append(ret, answer)
	}

	return
}
