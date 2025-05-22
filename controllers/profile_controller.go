package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile(c *gin.Context) {
	user := getCurrentUser(c)
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "email": user.Email})
}

func UpdateProfile(c *gin.Context) {
	user := getCurrentUser(c)
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	user.Username = input.Username
	user.Email = input.Email
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Обновлено"})
}

func ChangePassword(c *gin.Context) {
	user := getCurrentUser(c)
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный старый пароль"})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 10)
	user.Password = string(hashed)
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Пароль обновлён"})
}

func getCurrentUser(c *gin.Context) *models.User {
	tokenString, err := c.Cookie("token")
	if err != nil || tokenString == "" {
		return nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil || !token.Valid {
		return nil
	}

	claims := token.Claims.(jwt.MapClaims)
	var user models.User
	config.DB.First(&user, claims["user_id"])
	return &user
}

func DeleteProfile(c *gin.Context) {
	user := getCurrentUser(c)
	config.DB.Delete(&user)
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Аккаунт удалён"})
}

func UploadAvatar(c *gin.Context) {
	user := getCurrentUser(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет файла"})
		return
	}

	// Уникальное имя файла
	filename := fmt.Sprintf("avatar_%d_%s", user.ID, file.Filename)
	path := "static/avatars/" + filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки файла"})
		return
	}

	// Сохраняем путь в БД
	user.ImageURL = "/" + path
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Аватар загружен", "image_url": user.ImageURL})
}
