// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package service

import "github.com/b3log/routinepanic.com/model"

// Tag service.
var Tag = &tagService{}

type tagService struct {
}

func (srv *tagService) GetTopTags(size int) (ret []*model.Tag) {
	if err := db.Model(&model.Tag{}).Order("`question_count` DESC").Limit(size).Find(&ret).Error; nil != err {
		logger.Errorf("get top tags failed: %s", err)

		return
	}

	return
}

func (srv *tagService) GetTagByTitle(title string) *model.Tag {
	ret := &model.Tag{}
	if err := db.Where("`title` = ?", title).First(ret).Error; nil != err {
		return nil
	}

	return ret
}
