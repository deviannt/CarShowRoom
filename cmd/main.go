package main

import (
	"autosalon/config"
	"autosalon/models"
	"autosalon/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения по умолчанию")
	}

	// Подключение к базе данных
	config.ConnectDB()

	// Автоматическая миграция таблиц
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Car{},
	)
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	// Запуск маршрутизатора
	r := routes.SetupRouter()

	// Получаем порт из переменных окружения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(r.Run(":" + port))
}
