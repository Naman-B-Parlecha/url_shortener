package services

import "github.com/redis/go-redis/v9"

type MetricsService struct {
	client *redis.Client
}

func NewMetricsService(client *redis.Client) *MetricsService {
	return &MetricsService{
		client: client,
	}
}

// any service i make
