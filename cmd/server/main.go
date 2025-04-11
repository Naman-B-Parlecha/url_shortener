package main

import (
	"github/Naman-B-Parlecha/url-shortener/internal/routes"
	"github/Naman-B-Parlecha/url-shortener/pkg/redis"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	redisStore := redis.NewRedisStore()
	if err := redisStore.Ping(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r, nil)
	r.Run(":8080")
}
