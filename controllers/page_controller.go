package controllers

import (
	"autosalon/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// ✅ Получение информации из токена
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

// 🌐 HTML: Страница регистрации
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

// 🌐 HTML: Страница входа
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

// 🌐 HTML: Главная страница со списком автомобилей
func ShowCarsPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Автомобили",
		"Content":         "cars.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 HTML: Профиль пользователя
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

// 🌐 HTML: Админ-панель пользователей
func ShowAdminUsersPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Админ-панель пользователей",
		"Content":         "admin_users.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 HTML: Страница добавления авто
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

// 🌐 HTML: Редактирование авто
func ShowCarEditPage(c *gin.Context) {
	auth, username, role := getUserInfo(c)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"Title":           "Редактировать автомобиль",
		"Content":         "car_edit.html",
		"IsAuthenticated": auth,
		"Username":        username,
		"Role":            role,
	})
}

// 🌐 HTML: Список авто в админке
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
