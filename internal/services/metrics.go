package services

import (
	"strings"

	"github.com/redis/go-redis/v9"
)

type MetricsService struct {
	client *redis.Client
}

func NewMetricsService(client *redis.Client) *MetricsService {
	return &MetricsService{
		client: client,
	}
}

func (s *MetricsService) DomainCount() (map[string]int, error) {
	domains, err := s.client.Keys(ctx, "domain:*").Result()
	if err != nil {
		return nil, err
	}

	domainCounts := make(map[string]int)
	for _, domain := range domains {
		count, err := s.client.Get(ctx, domain).Int()
		if err != nil {
			return nil, err
		}
		domainName := strings.TrimPrefix(domain, "domain:")
		domainCounts[domainName] = count
	}

	return domainCounts, nil

}
