package config

// 🔐 Глобальный JWT-ключ
var JwtKey = []byte("secret_key") // Лучше заменить на os.Getenv("JWT_SECRET")
