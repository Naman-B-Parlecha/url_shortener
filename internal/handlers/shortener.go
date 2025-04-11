package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"
	"github/Naman-B-Parlecha/url-shortener/models"

	"github.com/gin-gonic/gin"
)

type ShortenerHandler struct {
	shortenerService *services.ShortenerService
}

func NewShortenerHandler(service *services.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		shortenerService: service,
	}
}

func (h *ShortenerHandler) ShortenURL(c *gin.Context) {
	var req models.URL

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	shortURL, err := h.shortenerService.ShortenURL(req.Url)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to shorten URL"})
		return
	}
	c.JSON(200, gin.H{"short_url": shortURL})
}
