// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

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
