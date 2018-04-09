package model

// Question model.
type Question struct {
	Model

	Title     string `gorm:"type:text" json:"title"`
	Tags      string `gorm:"type:text" json:"tags"`
	Content   string `gorm:"type:mediumtext" json:"content"`
	Path      string `gorm:"type:text" json:"path"`
	Source    int    `sql:"index" json:"source"`
	SourceID  string `gorm:"size:255" sql:"index"`
	SourceURL string `gorm:"size:255" sql:"index" json:"sourceURL"`
}

// Sources.
const (
	SourceStackOverflow = iota
)
