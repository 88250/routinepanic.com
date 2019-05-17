// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

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
