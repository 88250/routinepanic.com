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

// Word model.
type Word struct {
	Model

	Name    string `gorm:"size:128" sql:"unique_index" json:"name"`
	PhAm    string `gorm:"size:255" json:"phAm"`
	PhAmMp3 string `gorm:"size:255" json:"phAmMP3"`
	PhEn    string `gorm:"size:255" json:"phEn"`
	PhEnMp3 string `gorm:"size:255" json:"phEnMP3"`
	Means   string `gorm:"type:text" json:"means"`
}
