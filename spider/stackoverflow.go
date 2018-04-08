package spider

import (
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

	doc.Find("#question-header h1").Each(func(i int, s *goquery.Selection) {
		question.Title = s.Text()
	})
	tags := ""
	doc.Find(".post-taglist a").Each(func(i int, s *goquery.Selection) {
		tags += s.Text() + ","
	})
	if 0 < len(tags) {
		tags = tags[:len(tags)-1]
	}
	question.Tags = tags
	doc.Find("#question .post-text").Each(func(i int, s *goquery.Selection) {
		question.Content, _ = s.Html()
	})
	doc.Find("#answers .answer .post-text").EachWithBreak(func(i int, s *goquery.Selection) bool {
		content, _ := s.Html()
		answer := &model.Answer{
			Content: content,
			Source:  model.SourceStackOverflow,
		}
		answers = append(answers, answer)

		return i < 2
	})

	return
}
