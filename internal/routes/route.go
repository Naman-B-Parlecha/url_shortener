package routes

import (
	"github/Naman-B-Parlecha/url-shortener/internal/config"
	"github/Naman-B-Parlecha/url-shortener/internal/handlers"
	"github/Naman-B-Parlecha/url-shortener/internal/services"
	"github/Naman-B-Parlecha/url-shortener/pkg/redis"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config) {

	redisClient := redis.NewRedisStore()
	shortenerService := services.NewShortenerService(redisClient.Client)
	shortenerHandler := handlers.NewShortenerHandler(shortenerService)
	r.Group("/shorten")
	{
		r.POST("/", shortenerHandler.ShortenURL)
	}

	redirectService := services.NewRedirectService(redisClient.Client)
	redirectHandler := handlers.NewRedirectHandler(redirectService)
	r.Group("/:shortened")
	{
		r.GET("/", redirectHandler.Redirect)
	}

	metricsService := services.NewMetricsService(redisClient.Client)
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	r.Group("/metrics")
	{
		r.GET("/", func(c *gin.Context) {})
	}
}
