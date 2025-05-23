package middleware

import (
	"autosalon/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен: сначала из Authorization, потом из cookie
		authHeader := c.GetHeader("Authorization")
		var tokenString string

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Пробуем из cookie
			cookie, err := c.Cookie("token")
			if err != nil || cookie == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Нет токена"})
				return
			}
			tokenString = cookie
		}

		// Парсим токен
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return config.JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			return
		}

		// Проверяем claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Некорректные данные токена"})
			return
		}

		// Сохраняем userID и role в контекст
		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Set("userRole", claims["role"].(string))

		c.Next()
	}
}
