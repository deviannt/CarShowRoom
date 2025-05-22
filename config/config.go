package config

import (
	"autosalon/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	maxRetries := 10

	for i := 1; i <= maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Успешно подключено к базе данных")

			// 🔧 Автоматическая миграция таблиц
			err = DB.AutoMigrate(&models.User{}, &models.Car{})
			if err != nil {
				log.Fatalf("❌ Ошибка миграции: %v", err)
			}

			return
		}
		log.Printf("⏳ Попытка %d из %d: не удалось подключиться к БД: %v", i, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Не удалось подключиться к базе данных после нескольких попыток")
}
