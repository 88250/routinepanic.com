// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"html/template"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/model"
	"github.com/vinta/pangu"
)

type question struct {
	ID          uint64
	Path        string
	Title       string
	Description string
	Tags        []*tag
	Content     template.HTML
}

type tag struct {
	Title string
}

type answer struct {
	ID      uint64
	Content template.HTML
}

func questionsVos(qModels []*model.Question) (ret []*question) {
	for _, qModel := range qModels {
		q := questionVo(qModel)
		ret = append(ret, q)
	}

	return
}

func questionVo(qModel *model.Question) (ret *question) {
	desc := ""
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(qModel.ContentZhCN))
	if nil != err {
		logger.Errorf("extract description failed: " + err.Error())
	} else {
		desc = doc.Find("p").First().Text()
	}

	ret = &question{
		ID:          qModel.ID,
		Path:        qModel.Path,
		Title:       pangu.SpacingText(qModel.TitleZhCN),
		Description: desc,
		Content:     template.HTML(panguSpace(qModel.ContentZhCN)),
	}
	tagStrs := strings.Split(qModel.Tags, ",")
	for _, tagTitle := range tagStrs {
		ret.Tags = append(ret.Tags, &tag{Title: tagTitle})
	}

	return
}

func answersVos(aModels []*model.Answer) (ret []*answer) {
	for _, aModel := range aModels {
		a := answerVo(aModel)
		ret = append(ret, a)
	}

	return
}

func answerVo(aModel *model.Answer) (ret *answer) {
	ret = &answer{
		ID:      aModel.ID,
		Content: template.HTML(panguSpace(aModel.ContentZhCN)),
	}

	return
}

func panguSpace(html string) string {
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
	ret, err := doc.Find("body").Html()
	if nil != err {
		logger.Errorf("pangu space failed: " + err.Error())

		return html
	}

	return ret
}
