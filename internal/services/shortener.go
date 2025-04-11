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
	ctx    context.Context
}

func NewShortenerService(client *redis.Client, c context.Context) *ShortenerService {
	return &ShortenerService{
		client: client,
		ctx:    c,
	}
}

func (s *ShortenerService) ShortenURL(url string) (string, error) {
	shorturl, err := s.client.Get(s.ctx, url).Result()

	if err == redis.Nil {
		shorturl, err = generateShortenedUrl(url)
		if err != nil {
			return "", err
		}
		err = s.client.Set(s.ctx, shorturl, url, 0).Err()
		if err != nil {
			return "", err
		}
		domain, err := s.getDomain(url)
		if err != nil {
			return "", err
		}

		err = s.client.Incr(s.ctx, fmt.Sprintf("domain:%s", domain)).Err()
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
