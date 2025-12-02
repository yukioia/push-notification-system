package main

import (
	"log"
	"os"
	"push-server/src/internal/http"
	"push-server/src/internal/repo"
	"push-server/src/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к БД
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Репозитории
	subRepoPg := repo.NewPostgresSubscriptionRepo(db)
	subRepo := repo.NewSubscriptionRepoAdapter(subRepoPg)

	// Сервисы
	subService := service.NewSubscriptionService(subRepo)

	// HTTP хендлеры
	h := http.NewHandler(subService)
	r := gin.Default()

	// Статика (главная страница)
	r.StaticFile("/", "./web/index.html")
	r.Static("/static", "./web") // если есть дополнительные файлы, например JS/CSS

	// Роуты API
	h.RegisterRoutes(r)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
