package models

import (
	"time"
)

type Post struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	AuthorID    uint
	Author      User
	Approved    bool
	CreatedAt   time.Time
}
