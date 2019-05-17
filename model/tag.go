// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package model

// Tag model.
type Tag struct {
	Model

	Title         string `gorm:"size:128" json:"title"`
	QuestionCount int    `json:"questionCount"`
}
