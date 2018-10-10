// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import "github.com/b3log/routinepanic.com/model"

// User service.
var User = &userService{}

type userService struct {
}

func (srv *userService) Get(id uint64) (ret *model.User) {
	ret = &model.User{}
	db.Where("`id` = ?", id).First(ret)

	return
}

func (srv *userService) AddOrUpdate(user *model.User) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err = tx.Where("`name` = ?", user.Name).
		Assign(model.User{
			Name:   user.Name,
			Avatar: user.Avatar,
		}).FirstOrCreate(user).Error; nil != err {
		return
	}

	return nil
}
