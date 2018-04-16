// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"html/template"
	"strings"

	"github.com/b3log/routinepanic.com/model"
)

type question struct {
	ID      uint64
	Path    string
	Title   string
	Tags    []*tag
	Content template.HTML
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
	ret = &question{
		ID:      qModel.ID,
		Path:    qModel.Path,
		Title:   qModel.TitleZhCN,
		Content: template.HTML(qModel.ContentZhCN),
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
		Content: template.HTML(aModel.ContentZhCN),
	}

	return
}
