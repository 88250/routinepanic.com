// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package model

// Revision model.
type Revision struct {
	Model

	DataType    int     `json:"dataType"`
	DataID      uint64  `json:"dataID"`
	Data        string  `gorm:"type:mediumtext" json:"data"`
	AuthorID    uint64  `json:"authorID"`
	Distance    int     `json:"distance"`
	JaroWinkler float64 `json:"jaroWinkler"`
}
