package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model
	Brand       string  `json:"brand"`                           // Марка
	ModelName   string  `json:"model" gorm:"column:model"`       // Модель (в БД колонка называется 'model')
	Year        int     `json:"year"`                            // Год выпуска
	Price       float64 `json:"price"`                           // Цена
	Description string  `json:"description"`                     // Описание
	ImageURL    string  `json:"image_url"`                       // Ссылка на изображение
	UserID      uint    `json:"user_id"`                         // ID пользователя, добавившего машину
	Status      string  `json:"status" gorm:"default:'pending'"` // Статус модерации: pending, approved, rejected
}
