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

package model

// Question model.
type Question struct {
	Model

	TitleEnUS   string `gorm:"type:text" json:"titleEnUS"`
	TitleZhCN   string `gorm:"type:text" json:"titleZhCN"`
	Tags        string `gorm:"type:text" json:"tags"`
	Votes       int    `json:"votes"`
	Views       int    `json:"views"`
	ContentEnUS string `gorm:"type:mediumtext" json:"contentEnUS"`
	ContentZhCN string `gorm:"type:mediumtext" json:"contentZhCN"`
	Path        string `gorm:"type:text" json:"path"`
	Source      int    `sql:"index" json:"source"`
	SourceID    string `gorm:"size:255" sql:"index"`
	SourceURL   string `gorm:"size:255" sql:"index" json:"sourceURL"`
	AuthorName  string `json:"authorName"`
	AuthorURL   string `gorm:"size:255" json:"authorURL"`
}

// Sources.
const (
	SourceStackOverflow = iota
)

// Data types
const (
	DataTypeQuestion = iota
	DataTypeAnswer
)
