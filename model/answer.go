package model

// Answer model.
type Answer struct {
	Model

	Content   string `gorm:"type:mediumtext" json:"content"`
	Path      string `gorm:"type:text" json:"path"`
	Source    int    `sql:"index" json:"source"`
	SourceURL string `gorm:"type:text" json:"sourceURL"`
}
