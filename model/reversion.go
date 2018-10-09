// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package model

// Reversion model.
type Reversion struct {
	Model

	DataType   int    `json:"dataType"`
	DataId     uint64 `json:"dataId"`
	Data       string `gorm:"type:mediumtext" json:"data"`
	AuthorName string `json:"authorName"`
}
