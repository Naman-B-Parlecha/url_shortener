package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	metricsService *services.MetricsService
}

func NewMetricsHandler(service *services.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: service,
	}
}

func (h *MetricsHandler) GetMetrics(c *gin.Context) {}
