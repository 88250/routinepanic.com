// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"github.com/b3log/routinepanic.com/service"
	"html/template"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/util"
	"github.com/vinta/pangu"
)

type question struct {
	ID           uint64
	Path         string
	Title        string
	Description  string
	Tags         []*tag
	Content      template.HTML
	SourceURL    string
	Contributors []*contributor
}

type tag struct {
	Title string
}

type contributor struct {
	Name   string
	Avatar string
}

type answer struct {
	ID           uint64
	Content      template.HTML
	Contributors []*contributor
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

	contributorUsers := service.QnA.QContributors(qModel)
	ret.Contributors = []*contributor{}
	for _, contributorUser := range contributorUsers {
		ret.Contributors = append(ret.Contributors, &contributor{
			Name:   contributorUser.Name,
			Avatar: contributorUser.Avatar,
		})
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
