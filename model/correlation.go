// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

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
