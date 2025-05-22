package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username  string `json:"username" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"-"`
	Role      string `json:"role"`
	ImageURL  string `json:"image_url"`
	IsBlocked bool   `json:"is_blocked"`
}
