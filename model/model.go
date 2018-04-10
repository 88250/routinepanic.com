// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package model

import "time"

// Model represents meta data of entity.
type Model struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}
