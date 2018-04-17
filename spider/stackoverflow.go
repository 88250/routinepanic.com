// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package spider

import (
	"fmt"
	"os"
	"strconv"
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

func (s *stackOverflow) ParseQuestionsByVotes(fromPage, toPage int) (ret []*QnA) {
	for i := fromPage; i <= toPage; i++ {
		qnas := s.ParseQuestions(fmt.Sprintf("https://stackoverflow.com/questions?page=%d&sort=votes", i))
		if nil != qnas {
			ret = append(ret, qnas...)
		}

		logger.Infof("parsed voted questions [page=%d]", i)
	}

	return
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
		qna := s.ParseQuestion("https://stackoverflow.com" + url)
		if nil == qna {
			continue
		}

		ret = append(ret, qna)

		logger.Infof("parsed question #%d [%s]", i, qna.Question.TitleEnUS)
	}

	return ret
}

func (s *stackOverflow) ParseQuestion(url string) *QnA {
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

	question := &model.Question{}
	var answers []*model.Answer

	urlParts := strings.Split(url, "/")
	question = &model.Question{
		Path:      urlParts[len(urlParts)-1],
		Source:    model.SourceStackOverflow,
		SourceURL: url,
	}
	questionSrcID, _ := doc.Find("#question").Attr("data-questionid")
	question.SourceID = questionSrcID
	doc.Find("#question-header h1").Each(func(i int, s *goquery.Selection) {
		question.TitleEnUS = strings.TrimSpace(s.Text())
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
		question.ContentEnUS, _ = s.Html()
		question.ContentEnUS = strings.TrimSpace(question.ContentEnUS)
	})
	votesStr := doc.Find(".vote-count-post.high-scored-post").First().Text()
	question.Votes, err = strconv.Atoi(votesStr)
	if nil != err {
		logger.Errorf("parse [%s] failed: ", url, err)

		return nil
	}
	info := doc.Find("#qinfo").First().Text()
	viewsStr := strings.TrimSpace(between(info, "viewed", "times"))
	viewsStr = strings.Replace(viewsStr, ",", "", -1)
	question.Views, err = strconv.Atoi(viewsStr)
	if nil != err {
		logger.Errorf("parse [%s] failed: ", url, err)

		return nil
	}
	doc.Find("#answers .answer").EachWithBreak(func(i int, s *goquery.Selection) bool {
		answerSrcID, _ := s.Attr("data-answerid")
		votesStr := s.Find(".vote-count-post.high-scored-post").First().Text()
		votes, _ := strconv.Atoi(votesStr)
		content, _ := s.Find(".post-text").Html()
		content = strings.TrimSpace(content)
		answer := &model.Answer{
			Votes:       votes,
			ContentEnUS: content,
			Source:      model.SourceStackOverflow,
			SourceID:    answerSrcID,
		}
		answers = append(answers, answer)

		return i < 2
	})

	return &QnA{Question: question, Answers: answers}
}

func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
