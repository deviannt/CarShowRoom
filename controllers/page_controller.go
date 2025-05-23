package controllers

import (
	"autosalon/config"
	"autosalon/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// 🔐 Получение информации из JWT токена (из cookie)
func getUserInfo(c *gin.Context) (bool, string, string) {
	tokenString, err := c.Cookie("token")
	if err != nil || tokenString == "" {
		return false, "", ""
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})
	if err != nil || !token.Valid {
		return false, "", ""
	}

	claims := token.Claims.(jwt.MapClaims)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)

	return true, username, role
}

// 🌐 Регистрация
func ShowRegisterPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Регистрация",
		"Content":         "register.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Вход
func ShowLoginPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Вход",
		"Content":         "login.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Главная: только одобренные посты
func ShowCarsPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)

	var approvedPosts []models.Post
	config.DB.Preload("Author").Where("approved = true").Order("created_at desc").Find(&approvedPosts)

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Автомобили",
		"Content":         "cars.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
		"Posts":           approvedPosts,
	})
}

// 🌐 Профиль
func ShowProfilePage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Профиль",
		"Content":         "profile.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Страница добавления автомобиля
func ShowCarAddPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Добавить автомобиль",
		"Content":         "car_add.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Страница редактирования автомобиля
func ShowCarEditPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Редактировать авто",
		"Content":         "car_edit.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Админ: список машин
func ShowAdminCarsPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Админ: Автомобили",
		"Content":         "admin_cars.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Админ: список пользователей
func ShowAdminUsersPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Пользователи",
		"Content":         "admin_users.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Админ: модерация постов
func ShowAdminPostsPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Посты на модерации",
		"Content":         "admin_posts.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 Страница поддержки (общая для всех ролей)
func ShowSupportPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Поддержка",
		"Content":         "support.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}
