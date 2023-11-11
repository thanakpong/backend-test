package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	Id         int       `json:"Id"`
	Title      string    `json:"Title"`
	Content    string    `json:"content"`
	Published  bool      `json:"Published"`
	View_count int       `json:"View_count"`
	Created_at time.Time `json:"Created_at"`
	Update_at  time.Time `json:"Update_at"`
}

func MigratePost(db *gorm.DB) error {
	err := db.AutoMigrate(&Post{})
	return err
}
