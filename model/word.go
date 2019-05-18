// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

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
