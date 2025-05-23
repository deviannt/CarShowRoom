package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 📄 Получить всех пользователей (admin)
func ListUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

// 🚫 Заблокировать пользователя (admin)
func BlockUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if user.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нельзя заблокировать супер-админа"})
		return
	}
	user.IsBlocked = true
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь заблокирован"})
}

// 🔓 Разблокировать пользователя (superadmin)
func UnblockUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if !user.IsBlocked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь не заблокирован"})
		return
	}
	user.IsBlocked = false
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь разблокирован"})
}

// ✏️ Изменить имя пользователя (admin)
func UpdateUsername(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректное имя"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	if user.Role == "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нельзя менять имя супер-админа"})
		return
	}

	var exists models.User
	if err := config.DB.Where("username = ? AND id != ?", input.Username, id).First(&exists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя уже занято"})
		return
	}

	config.DB.Model(&user).Update("username", input.Username)
	c.JSON(http.StatusOK, gin.H{"message": "Имя пользователя обновлено"})
}

// 🔄 Изменить роль пользователя (superadmin)
func SetUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	id := c.Param("id")
	if err := config.DB.Model(&models.User{}).Where("id = ?", id).Update("role", input.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления роли"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Роль обновлена"})
}

// ❌ Удалить пользователя (superadmin)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	config.DB.Delete(&models.User{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
}
