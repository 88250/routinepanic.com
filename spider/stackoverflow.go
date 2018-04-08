package spider

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/log"
	"github.com/corpix/uarand"
	"github.com/parnurzeal/gorequest"
)

var StackOverflow = &stackOverflow{}

type stackOverflow struct{}

type Question struct {
	Content string
}

// Logger
var logger = log.NewLogger(os.Stdout)

func (s *stackOverflow) ParseQuestion(url string) *Question {
	request := gorequest.New()
	response, body, errs := request.Set("User-Agent", uarand.GetRandom()).Get(url).End()
	if nil != errs {
		logger.Errorf("gets [%s] failed: %s", url, errs)

		return nil
	}
	if 200 != response.StatusCode {
		logger.Errorf("gets [%s] status code is [%d]", response.StatusCode)

		return nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if nil != err {
		logger.Errorf("parses [%s] failed: ", url, err)

		return nil
	}

	ret := &Question{}
	doc.Find("#question .post-text").Each(func(i int, s *goquery.Selection) {
		ret.Content, _ = s.Html()
	})

	return ret
}
