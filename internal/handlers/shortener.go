package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"

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

func (h *ShortenerHandler) ShortenURL(c *gin.Context) {}
