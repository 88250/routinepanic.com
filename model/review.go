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

import "time"

// Review model.
type Review struct {
	Model

	RevisionID uint64    `json:"revisionID"`
	Status     int       `json:"status"`
	ReviewerID uint64    `json:"reviewerID"`
	Memo       string    `gorm:"type:mediumtext" json:"memo"`
	ReviewedAt time.Time `json:"reviewedAt"`
}

const (
	ReviewStatusWaiting = iota
	ReviewStatusPassed
	ReviewStatusRejected
)
