package redis

import "github.com/redis/go-redis/v9"

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore() *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisStore{
		Client: client,
	}
}
