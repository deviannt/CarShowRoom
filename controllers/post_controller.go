package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// ✅ Создание поста (пользователь)
func CreatePost(c *gin.Context) {
	user := getCurrentUser(c)

	title := c.PostForm("title")
	description := c.PostForm("description")
	file, err := c.FormFile("image")
	if title == "" || description == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Все поля обязательны"})
		return
	}

	filename := fmt.Sprintf("post_%d_%s", user.ID, filepath.Base(file.Filename))
	path := filepath.Join("static/posts", filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки изображения"})
		return
	}

	post := models.Post{
		Title:       title,
		Description: description,
		AuthorID:    user.ID,
		ImageURL:    "/" + path,
		Approved:    false,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания поста"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост отправлен на модерацию"})
}

// ✅ Список неподтверждённых постов (для модерации)
func ListUnapprovedPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("Author").Where("approved = false").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// ✅ Подтверждение поста (админ)
func ApprovePost(c *gin.Context) {
	var post models.Post
	if err := config.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}
	post.Approved = true
	config.DB.Save(&post)
	c.JSON(http.StatusOK, gin.H{"message": "Пост одобрен"})
}

// ✅ Удаление поста (админ)
func DeletePost(c *gin.Context) {
	if err := config.DB.Delete(&models.Post{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления поста"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пост удалён"})
}

// ✏️ Обновление поста (админ)
func UpdatePost(c *gin.Context) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.Title == "" || input.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пост не найден"})
		return
	}

	post.Title = input.Title
	post.Description = input.Description
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления поста"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пост обновлён"})
}
