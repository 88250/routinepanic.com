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

import (
	"github.com/b3log/routinepanic.com/model"
	"github.com/jinzhu/gorm"
)

// Dictionary service.
var Dic = &dicService{}

type dicService struct {
}

func (srv *dicService) AddWord(word *model.Word) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err = tx.Where("`name` = ?", word.Name).
		Assign(model.Word{
			Name:    word.Name,
			PhAm:    word.PhAm,
			PhAmMp3: word.PhEnMp3,
			PhEn:    word.PhEn,
			PhEnMp3: word.PhEnMp3,
			Means:   word.Means,
		}).FirstOrCreate(word).Error; nil != err {
		return
	}

	return nil
}

func (srv *dicService) GetWord(name string) (ret *model.Word) {
	ret = &model.Word{}
	if err := db.Model(&model.Word{}).Where("`name` = ?", name).First(ret).Error; nil != err {
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("get word [%s] failed: "+err.Error(), name)
		}

		return nil
	}

	return
}
