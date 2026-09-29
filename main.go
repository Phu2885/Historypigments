package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rip/internal/app/handler"
	"rip/internal/app/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 1. Подключение к БД PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2. Инициализация Репозитория
	repo := repository.New(db)

	// Инициализация MinIO SDK клиента (ЛР3)
	minioHost := os.Getenv("MINIO_ENDPOINT")
	if minioHost == "" {
		minioHost = "localhost:9000"
	}
	minioUser := os.Getenv("MINIO_ROOT_USER")
	if minioUser == "" {
		minioUser = "root"
	}
	minioPass := os.Getenv("MINIO_ROOT_PASSWORD")
	if minioPass == "" {
		minioPass = "rootpassword"
	}

	minioClient, errMinio := minio.New(minioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPass, ""),
		Secure: false,
	})
	if errMinio == nil {
		repo.SetMinio(minioClient, "historypigments")
		log.Println("MinIO SDK Client initialized successfully")
	} else {
		log.Println("MinIO SDK initialization warning:", errMinio)
	}

	h := handler.New(repo)

	// 3. Роутинг
	mux := http.NewServeMux()

	// Статика и прокси к MinIO
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets"))))

	minioURL, _ := url.Parse("http://" + minioHost)
	proxy := httputil.NewSingleHostReverseProxy(minioURL)
	mux.Handle("/minio/", http.StripPrefix("/minio/", proxy))

	// ==========================================
	// REST API РУТИНГ (10 Доменных Эндпоинтов ЛР3)
	// ==========================================

	// /api/pigments
	mux.HandleFunc("/api/pigments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetPigmentsListAPI(w, r)
		case http.MethodPost:
			h.CreateDraftAPI(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /api/pigments/...
	mux.HandleFunc("/api/pigments/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/api/pigments/feed" {
			if r.Method == http.MethodGet {
				h.GetFeedAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		if path == "/api/pigments/draft" {
			if r.Method == http.MethodGet {
				h.GetDraftAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		if strings.HasSuffix(path, "/publish") {
			if r.Method == http.MethodPut || r.Method == http.MethodPost {
				h.PublishPigmentAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		if strings.HasSuffix(path, "/like") {
			if r.Method == http.MethodPost || r.Method == http.MethodPut {
				h.ToggleLikeAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /api/pigments/{id}
		switch r.Method {
		case http.MethodGet:
			h.GetPigmentByIDAPI(w, r)
		case http.MethodDelete:
			h.DeletePigmentAPI(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /api/users/...
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/api/users/register" {
			if r.Method == http.MethodPost {
				h.RegisterUserAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		if path == "/api/users/login" {
			if r.Method == http.MethodPost {
				h.LoginUserAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}
		if path == "/api/users/logout" {
			if r.Method == http.MethodPost {
				h.LogoutUserAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if strings.HasSuffix(path, "/pigments") {
			if r.Method == http.MethodGet {
				h.GetUserPigmentsAPI(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})

	fmt.Printf("Alchemy of Art History REST API serving on http://localhost:8080\n")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
