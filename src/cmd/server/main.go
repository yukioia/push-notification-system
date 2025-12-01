package main

import (
	"log"
	"os"
	"push-server/src/internal/http"
	"push-server/src/internal/repo"
	"push-server/src/internal/service"
	"push-server/src/internal/sse"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}

	repoPg, err := repo.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	sseMgr := sse.NewManager()
	svc := service.NewNotificationService(repoPg, sseMgr)
	h := http.NewHandler(svc, sseMgr)

	r := gin.Default()
	r.StaticFile("/", "./web/index.html")
	h.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
