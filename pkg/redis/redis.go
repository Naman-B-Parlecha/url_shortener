package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStore() *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisStore{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (s *RedisStore) Ping() error {
	_, err := s.Client.Ping(s.Ctx).Result()
	return err
}
