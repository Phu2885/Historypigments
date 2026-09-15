package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rip/internal/app/handler"
	"rip/internal/app/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 1. Подключение к БД
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. Инициализация Репозитория и Хендлера
	tmpl := template.Must(template.ParseGlob(filepath.Join("./templates", "*.html")))
	repo := repository.New(db)
	h := handler.New(repo, tmpl)

	// 3. Роутинг
	mux := http.NewServeMux()

	// Статика
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets"))))

	// GET методы
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/pigments-database-reconstruction", http.StatusFound)
	})
	mux.HandleFunc("/pigments-database-reconstruction", h.GetFeed)
	mux.HandleFunc("/pigment-database-analysis", h.GetAnalysis)
	mux.HandleFunc("/pigments-database", h.GetDatabase)

	// POST методы
	mux.HandleFunc("/create-draft", h.CreateDraft)
	mux.HandleFunc("/publish-pigment", h.Publish)
	mux.HandleFunc("/delete-pigment", h.Delete)

	fmt.Printf("Alchemy of Art History serving on http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
