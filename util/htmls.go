// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package util

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vinta/pangu"
)

func TuneHTML(html string) string {
	if "" == html {
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if nil != err {
		logger.Errorf("pangu space failed: " + err.Error())

		return html
	}
	doc.Find("*").Contents().FilterFunction(func(i int, ele *goquery.Selection) bool {
		if "#text" != goquery.NodeName(ele) {
			return false
		}
		parent := goquery.NodeName(ele.Parent())

		return parent != "code" && parent != "pre"
	}).Each(func(i int, ele *goquery.Selection) {
		text := pangu.SpacingText(ele.Text())
		text = pangu.SpacingText(text)
		ele.ReplaceWithHtml(text)
	})
	doc.Find("pre").Each(func(i int, ele *goquery.Selection) {
		html, _ := ele.Html()
		html = strings.TrimSpace(html)
		ele.SetHtml(html)
	})
	doc.Find("code").Each(func(i int, ele *goquery.Selection) {
		html, _ := ele.Html()
		html = strings.TrimSpace(html)
		html = strings.Replace(html, "<", "&lt;", -1)
		html = strings.Replace(html, ">", "&gt;", -1)
		ele.SetHtml(html)
	})
	ret, err := doc.Find("body").Html()
	if nil != err {
		logger.Errorf("pangu space failed: " + err.Error())

		return html
	}

	return ret
}
