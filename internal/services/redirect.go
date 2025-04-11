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
