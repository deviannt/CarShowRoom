package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ✅ Получить все сообщения (админ)
func GetSupportMessages(c *gin.Context) {
	var messages []models.SupportMessage
	if err := config.DB.Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить сообщения"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// ✅ Отправить сообщение (пользователь или админ)
func SendSupportMessage(c *gin.Context) {
	var input struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Введите сообщение"})
		return
	}

	username := c.GetString("username")
	role := c.GetString("role")

	msg := models.SupportMessage{
		Username: username,
		Role:     role,
		Message:  input.Message,
	}

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Сообщение отправлено"})
}
