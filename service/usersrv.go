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

// User service.
var User = &userService{}

type userService struct {
}

func (srv *userService) GetByName(name string) (ret *model.User) {
	ret = &model.User{}
	db.Where("`name` = ?", name).First(ret)

	return
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
			Name:     user.Name,
			Avatar:   user.Avatar,
			GithubId: user.GithubId,
		}).FirstOrCreate(user).Error; nil != err {
		return
	}

	return nil
}
