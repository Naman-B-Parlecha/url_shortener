package services

import "github.com/redis/go-redis/v9"

type ShortenerService struct {
	client *redis.Client
}

func NewShortenerService(client *redis.Client) *ShortenerService {
	return &ShortenerService{
		client: client,
	}
}
