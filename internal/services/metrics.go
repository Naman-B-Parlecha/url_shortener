package services

import (
	"context"
	"strings"

	"github.com/redis/go-redis/v9"
)

type MetricsService struct {
	client *redis.Client
	ctx    context.Context
}

func NewMetricsService(client *redis.Client, c context.Context) *MetricsService {
	return &MetricsService{
		client: client,
		ctx:    c,
	}
}

func (s *MetricsService) DomainCount() (map[string]int, error) {
	domains, err := s.client.Keys(s.ctx, "domain:*").Result()
	if err != nil {
		return nil, err
	}

	domainCounts := make(map[string]int)
	for _, domain := range domains {
		count, err := s.client.Get(s.ctx, domain).Int()
		if err != nil {
			return nil, err
		}
		domainName := strings.TrimPrefix(domain, "domain:")
		domainCounts[domainName] = count
	}

	return domainCounts, nil

}
