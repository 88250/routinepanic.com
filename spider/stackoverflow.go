// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package spider

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/model"
	"github.com/corpix/uarand"
	"github.com/parnurzeal/gorequest"
)

var StackOverflow = &stackOverflow{}

type stackOverflow struct{}

// Logger
var logger = log.NewLogger(os.Stdout)

// QnA represents a question and its answers.
type QnA struct {
	Question *model.Question
	Answers  []*model.Answer
}

func (s *stackOverflow) ParseQuestionsByVotes(fromPage, toPage int) []*QnA {
	var ret []*QnA
	for i := fromPage; i <= toPage; i++ {
		qnas := s.ParseQuestions(fmt.Sprintf("https://stackoverflow.com/questions?page=%d&sort=votes", i))
		if nil != qnas {
			ret = append(ret, qnas...)
		}

		logger.Infof("parsed voted questions [page=%d]", i)
	}

	return ret
}

func (s *stackOverflow) ParseQuestions(url string) []*QnA {
	request := gorequest.New()
	response, body, errs := request.Set("User-Agent", uarand.GetRandom()).Get(url).End()
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d]", response.StatusCode)

		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if nil != err {
		logger.Errorf("parse [%s] failed: ", url, err)

		return nil
	}

	var questionURLs []string
	doc.Find("#questions .summary h3 a").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		questionURLs = append(questionURLs, url)
	})

	var ret []*QnA
	for i, url := range questionURLs {
		question, answers := s.ParseQuestion("https://stackoverflow.com" + url)
		qna := &QnA{Question: question, Answers: answers}
		ret = append(ret, qna)

		logger.Infof("parsed question #%d [%s]", i, question.Title)
	}

	return ret
}

func (s *stackOverflow) ParseQuestion(url string) (question *model.Question, answers []*model.Answer) {
	request := gorequest.New()
	response, body, errs := request.Set("User-Agent", uarand.GetRandom()).Get(url).End()
	if nil != errs {
		logger.Errorf("get [%s] failed: %s", url, errs)

		return nil, nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("get [%s] status code is [%d]", response.StatusCode)

		return nil, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if nil != err {
		logger.Errorf("parse [%s] failed: ", url, err)

		return nil, nil
	}

	question = &model.Question{
		Source:    model.SourceStackOverflow,
		SourceURL: url,
	}
	questionSrcID, _ := doc.Find("#question").Attr("data-questionid")
	question.SourceID = questionSrcID
	doc.Find("#question-header h1").Each(func(i int, s *goquery.Selection) {
		question.Title = strings.TrimSpace(s.Text())
	})
	tags := ""
	doc.Find(".post-taglist a").Each(func(i int, s *goquery.Selection) {
		tags += strings.TrimSpace(s.Text()) + ","
	})
	if 0 < len(tags) {
		tags = tags[:len(tags)-1]
	}
	question.Tags = tags
	doc.Find("#question .post-text").Each(func(i int, s *goquery.Selection) {
		question.Content, _ = s.Html()
		question.Content = strings.TrimSpace(question.Content)
	})
	doc.Find("#answers .answer").EachWithBreak(func(i int, s *goquery.Selection) bool {
		answerSrcID, _ := s.Attr("data-answerid")
		content, _ := s.Find(".post-text").Html()
		content = strings.TrimSpace(content)
		answer := &model.Answer{
			Content:  content,
			Source:   model.SourceStackOverflow,
			SourceID: answerSrcID,
		}
		answers = append(answers, answer)

		return i < 2
	})

	return
}
