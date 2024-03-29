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

// Correlation types.
const (
	CorrelationQuestionTag = iota
)

// Correlation model.
//   id1(question_id) - id2(tag_id)
type Correlation struct {
	Model

	ID1  uint64 `json:"id1"`
	ID2  uint64 `json:"id2"`
	Str1 string `gorm:"size:255" json:"str1"`
	Str2 string `gorm:"size:255" json:"str2"`
	Str3 string `gorm:"size:255" json:"str3"`
	Str4 string `gorm:"size:255" json:"str4"`
	Int1 int    `json:"int1"`
	Int2 int    `json:"int2"`
	Int3 int    `json:"int3"`
	Int4 int    `json:"int4"`
	Type int    `json:"type"`
}
