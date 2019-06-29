// 协慌网 - 专注编程问答汉化 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
