package main

import (
	"github/Naman-B-Parlecha/url-shortener/internal/routes"
	"github/Naman-B-Parlecha/url-shortener/pkg/redis"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(time.Second), 5)

func rateLimitMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests. Please slow down!",
		})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	r := gin.Default()

	redisStore := redis.NewRedisStore()
	if err := redisStore.Ping(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(rateLimitMiddleware)
	routes.SetupRoutes(r, nil)
	r.Run(":8080")
}
