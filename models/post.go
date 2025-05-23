package models

import (
	"time"
)

type Post struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	AuthorID    uint      `json:"author_id"`
	Author      User      `json:"author" gorm:"foreignKey:AuthorID"`
	Approved    bool      `json:"approved"`
	CreatedAt   time.Time `json:"created_at"`
}
