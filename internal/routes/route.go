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
	shortenerService := services.NewShortenerService(redisClient.Client, redisClient.Ctx)
	shortenerHandler := handlers.NewShortenerHandler(shortenerService)
	shortenRoute := r.Group("/shorten")
	{
		shortenRoute.POST("/", shortenerHandler.ShortenURL)
	}

	redirectService := services.NewRedirectService(redisClient.Client, redisClient.Ctx)
	redirectHandler := handlers.NewRedirectHandler(redirectService)
	redirectRoute := r.Group("/shortened")
	{
		redirectRoute.GET("/:id", redirectHandler.Redirect)
	}

	metricsService := services.NewMetricsService(redisClient.Client, redisClient.Ctx)
	metricsHandler := handlers.NewMetricsHandler(metricsService)
	metricsRotues := r.Group("/metrics")
	{
		metricsRotues.GET("/", metricsHandler.GetMetrics)
	}
}
