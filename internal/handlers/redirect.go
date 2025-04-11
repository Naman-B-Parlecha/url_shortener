package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"

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

func (h *RedirectHandler) Redirect(c *gin.Context) {}
