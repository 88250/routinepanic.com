package model

// Answer model.
type Answer struct {
	Model

	QuestionID uint64 `sql:"index" json:"questionID"`
	Content    string `gorm:"type:mediumtext" json:"content"`
	Path       string `gorm:"type:text" json:"path"`
	Source     int    `sql:"index" json:"source"`
	SourceID   string `gorm:"size:255" sql:"index"`
	SourceURL  string `gorm:"size:255" sql:"index" json:"sourceURL"`
}
