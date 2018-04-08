package spider

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/log"
	"os"
)

var StackOverflow = &stackOverflow{}

type stackOverflow struct{}

type Question struct {
	Content string
}

// Logger
var logger = log.NewLogger(os.Stdout)

func (s *stackOverflow) ParseQuestion(url string) *Question {
	res, err := http.Get(url)
	if nil != err {
		logger.Errorf("gets [%s] failed: %s", url, err)

		return nil
	}
	defer res.Body.Close()

	if 200 != res.StatusCode {
		logger.Errorf("gets [%s] status code is [%d]", res.StatusCode)

		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
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
