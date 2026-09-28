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

	// Добавление дефолтного пользователя-создателя
	var defaultUser ds.User
	if err := db.Where("id = ? OR login = ?", 1, "test").First(&defaultUser).Error; err == gorm.ErrRecordNotFound {
		defaultUser = ds.User{
			Name:     "Эксперт Искусствовед",
			Login:    "test",
			Password: "test123",
		}
		_ = db.Create(&defaultUser).Error
	}

	fmt.Println("Database migration completed successfully")
}
