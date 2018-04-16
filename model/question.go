// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package model

// Question model.
type Question struct {
	Model

	TitleEnUS   string `gorm:"type:text" json:"titleEnUS"`
	TitleZhCN   string `gorm:"type:text" json:"titleZhCN"`
	Tags        string `gorm:"type:text" json:"tags"`
	ContentEnUS string `gorm:"type:mediumtext" json:"contentEnUS"`
	ContentZhCN string `gorm:"type:mediumtext" json:"contentZhCN"`
	Path        string `gorm:"type:text" json:"path"`
	Source      int    `sql:"index" json:"source"`
	SourceID    string `gorm:"size:255" sql:"index"`
	SourceURL   string `gorm:"size:255" sql:"index" json:"sourceURL"`
}

// Sources.
const (
	SourceStackOverflow = iota
)
