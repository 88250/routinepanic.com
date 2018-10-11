// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package model

// Review model.
type Review struct {
	Model

	RevisionID uint64 `json:"revisionID"`
	ReviewerID uint64 `json:"reviewerID"`
	Memo       string `gorm:"type:mediumtext" json:"memo"`
	Status     int    `json:"status"`
}

const (
	ReviewStatusWaiting = iota
	ReviewStatusPassed
	ReviewStatusRejected
)
