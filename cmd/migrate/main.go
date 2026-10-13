package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rip/internal/app/ds"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Миграция схем
	err = db.AutoMigrate(&ds.User{}, &ds.Pigment{}, &ds.Like{})
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	// Добавление / обновление дефолтного пользователя (test / test123)
	var testUser ds.User
	if err := db.Where("login = ?", "test").First(&testUser).Error; err == gorm.ErrRecordNotFound {
		testUser = ds.User{
			Name:     "Эксперт Искусствовед",
			Login:    "test",
			Password: "test123",
			Role:     "admin",
		}
		db.Create(&testUser)
	} else {
		db.Model(&testUser).Updates(map[string]interface{}{
			"password": "test123",
			"role":     "admin",
		})
	}

	fmt.Println("Database migration completed successfully")
}
