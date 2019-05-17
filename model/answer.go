// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package model

// Answer model.
type Answer struct {
	Model

	QuestionID  uint64 `sql:"index" json:"questionID"`
	Votes       int    `json:"votes"`
	ContentEnUS string `gorm:"type:mediumtext" json:"contentEnUS"`
	ContentZhCN string `gorm:"type:mediumtext" json:"contentZhCN"`
	Path        string `gorm:"type:text" json:"path"`
	Source      int    `sql:"index" json:"source"`
	SourceID    string `gorm:"size:255" sql:"index"`
	SourceURL   string `gorm:"size:255" sql:"index" json:"sourceURL"`
	AuthorName  string `json:"authorName"`
	AuthorURL   string `gorm:"size:255" json:"authorURL"`
}
