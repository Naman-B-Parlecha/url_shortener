package services

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/url"
	"strings"

	"github.com/redis/go-redis/v9"
)

type ShortenerService struct {
	client *redis.Client
}

func NewShortenerService(client *redis.Client) *ShortenerService {
	return &ShortenerService{
		client: client,
	}
}

var (
	ctx = context.Background()
)

func (s *ShortenerService) ShortenURL(url string) (string, error) {
	shorturl, err := s.client.Get(ctx, url).Result()

	if err == redis.Nil {
		shorturl, err = generateShortenedUrl(url)
		if err != nil {
			return "", err
		}
		err = s.client.Set(ctx, shorturl, url, 0).Err()
		if err != nil {
			return "", err
		}
		domain, err := s.getDomain(url)
		if err != nil {
			return "", err
		}

		err = s.client.Incr(ctx, fmt.Sprintf("domain:%s", domain)).Err()
		if err == redis.Nil {
			return "", fmt.Errorf("domain not found")
		}
	} else if err != nil {
		return "", err
	}
	return shorturl, nil
}

func generateShortenedUrl(url string) (string, error) {
	h := fnv.New32a()
	h.Write([]byte(url))

	return fmt.Sprintf("%x", h.Sum32()), nil
}

func (s *ShortenerService) getDomain(originalURL string) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(parsedURL.Host, "www."), nil
}
