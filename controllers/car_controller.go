package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// получить список авто с фильтрацией
func GetCars(c *gin.Context) {
	var cars []models.Car
	query := c.Query("q")
	max := c.Query("max")

	db := config.DB

	if query != "" {
		db = db.Where("brand ILIKE ? OR model ILIKE ?", "%"+query+"%", "%"+query+"%")
	}
	if max != "" {
		db = db.Where("price <= ?", max)
	}

	db.Find(&cars)
	c.JSON(http.StatusOK, cars)
}

// Получить одно авто
func GetCar(c *gin.Context) {
	id := c.Param("id")
	var car models.Car
	if err := config.DB.First(&car, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Авто не найдено"})
		return
	}
	c.JSON(http.StatusOK, car)
}

// Создать авто
func CreateCar(c *gin.Context) {
	var car models.Car
	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	config.DB.Create(&car)
	c.JSON(http.StatusOK, car)
}

// Обновить авто
func UpdateCar(c *gin.Context) {
	id := c.Param("id")
	var car models.Car
	if err := config.DB.First(&car, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Авто не найдено"})
		return
	}

	var input models.Car
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	car.Brand = input.Brand
	car.Model = input.Model
	car.Year = input.Year
	car.Price = input.Price
	car.Description = input.Description
	car.ImageURL = input.ImageURL

	config.DB.Save(&car)
	c.JSON(http.StatusOK, car)
}

// Удалить авто
func DeleteCar(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Car{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Авто удалено"})
}
