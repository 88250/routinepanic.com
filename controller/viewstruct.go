// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"html/template"
	"strings"

	"github.com/b3log/routinepanic.com/service"

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
	Name           string
	Avatar         string
	ContriCount    int
	ContriDistance int
}

type answer struct {
	ID           uint64
	Content      template.HTML
	Contributors []*contributor
}

const QINIU_IMG_STYLE_AVATAR = "imageView2/1/w/64/h/64/interlace/0/q/100"

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

	revisions := service.QnA.QRevisions(qModel)
	contributorMap := map[uint64]*contributor{}

	for _, revision := range revisions {
		contributorId := revision.AuthorID
		if val, ok := contributorMap[contributorId]; !ok {
			contributor := &contributor{}
			user := service.User.Get(revision.AuthorID)
			contributor.Name = user.Name
			contributor.Avatar = user.Avatar
			contributor.ContriCount = 1
			contributor.ContriDistance = revision.Distance
			contributorMap[contributorId] = contributor
		} else {
			val.ContriCount += 1
			val.ContriDistance += revision.Distance
		}
	}

	for _, val := range contributorMap {
		ret.Contributors = append(ret.Contributors, val)
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

	contributorUsers := service.QnA.AContributors(aModel)
	ret.Contributors = []*contributor{}
	for _, contributorUser := range contributorUsers {
		ret.Contributors = append(ret.Contributors, &contributor{
			Name:   contributorUser.Name,
			Avatar: contributorUser.Avatar + "?" + QINIU_IMG_STYLE_AVATAR,
		})
	}

	return
}
