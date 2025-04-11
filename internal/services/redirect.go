package services

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedirectService struct {
	client *redis.Client

	ctx context.Context
}

func NewRedirectService(client *redis.Client, c context.Context) *RedirectService {
	return &RedirectService{
		client: client,
		ctx:    c,
	}
}

func (s *RedirectService) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := s.client.Get(s.ctx, shortURL).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return originalURL, nil
}
