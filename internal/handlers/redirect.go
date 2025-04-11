package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedirectHandler struct {
	redirectService *services.RedirectService
}

func NewRedirectHandler(service *services.RedirectService) *RedirectHandler {
	return &RedirectHandler{
		redirectService: service,
	}
}

func (h *RedirectHandler) Redirect(c *gin.Context) {
	shortenedID := c.Param("id")

	if shortenedID == "" {
		c.JSON(400, gin.H{"error": "Missing URL identifier"})
		return
	}

	originalURL, err := h.redirectService.GetOriginalURL(shortenedID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get original URL"})
		return
	}

	if originalURL == "" {
		c.JSON(404, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
