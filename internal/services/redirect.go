package services

import "github.com/redis/go-redis/v9"

type RedirectService struct {
	client *redis.Client
}

func NewRedirectService(client *redis.Client) *RedirectService {
	return &RedirectService{
		client: client,
	}
}

func (s *RedirectService) GetOriginalURL(shortURL string) (string, error) {
	originalURL, err := s.client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return originalURL, nil
}
