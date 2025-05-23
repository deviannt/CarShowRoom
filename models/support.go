package models

import "gorm.io/gorm"

type SupportMessage struct {
	gorm.Model
	Username string `json:"username"` // имя отправителя
	Role     string `json:"role"`     // роль (user, admin)
	Message  string `json:"message"`  // текст сообщения
}
