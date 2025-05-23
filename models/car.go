package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	Brand       string  `json:"brand"`
	ModelName   string  `json:"model" gorm:"column:model"`
	Year        int     `json:"year"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	UserID      uint    `json:"user_id"`                         // кто подал
	Status      string  `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
}
