package handlers

import (
	"github/Naman-B-Parlecha/url-shortener/internal/services"
	"sort"

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

func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	domainCount, err := h.metricsService.DomainCount()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get metrics"})
		return
	}

	type kv struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}

	var sortedKvs []kv
	for k, v := range domainCount {
		sortedKvs = append(sortedKvs, kv{k, v})
	}

	sort.Slice(sortedKvs, func(i, j int) bool {
		return sortedKvs[i].Value > sortedKvs[j].Value
	})

	// Limit to top 10
	var topKvs []kv
	for i, v := range sortedKvs {
		if i >= 10 {
			break
		}
		topKvs = append(topKvs, v)
	}
	c.JSON(200, gin.H{"metrics": topKvs})

}
