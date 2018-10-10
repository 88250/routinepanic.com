// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"html/template"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/util"
	"github.com/vinta/pangu"
)

type question struct {
	ID          uint64
	Path        string
	Title       string
	Description string
	Tags        []*tag
	Content     template.HTML
	SourceURL   string
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

	content := qModel.ContentZhCN
	content = util.TuneHTML(content)
	title := pangu.SpacingText(qModel.TitleZhCN)

	ret = &question{
		ID:          qModel.ID,
		Path:        qModel.Path,
		Title:       title,
		Description: desc,
		Content:     template.HTML(content),
		SourceURL:   qModel.SourceURL,
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
	content := aModel.ContentZhCN
	content = util.TuneHTML(content)

	ret = &answer{
		ID:      aModel.ID,
		Content: template.HTML(content),
	}

	return
}
