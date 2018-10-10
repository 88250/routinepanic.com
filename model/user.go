// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package model

// User model.
type User struct {
	Model

	Name   string `gorm:"size:32" json:"name"`
	Avatar string `gorm:"size:255" json:"avatar"`
}
