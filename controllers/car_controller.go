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
	var input struct {
		Brand       string  `json:"brand"`
		ModelName   string  `json:"model"`
		Year        int     `json:"year"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		ImageURL    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	userID := c.GetUint("userID")

	car := models.Car{
		Brand:       input.Brand,
		ModelName:   input.ModelName,
		Year:        input.Year,
		Price:       input.Price,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		UserID:      userID,
		Status:      "pending",
	}

	if err := config.DB.Create(&car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Объявление отправлено на модерацию"})
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

func GetMyCars(c *gin.Context) {
	userID := c.GetUint("userID")

	var cars []models.Car
	if err := config.DB.Where("user_id = ?", userID).Find(&cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении машин"})
		return
	}

	c.JSON(http.StatusOK, cars)
}
