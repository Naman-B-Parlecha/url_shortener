package util

import (
	"crypto/md5"
	"math/rand"
	"time"
)

const base62Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func HashToBase62(longurl string, length int) string {
	hash := md5.Sum([]byte(longurl))

	var num uint64
	for i := 0; i < 8 && i < len(hash); i++ {
		num = (num << 8) | uint64(hash[i])
	}

	base62 := make([]byte, 0, length)
	for num > 0 && len(base62) < length {
		remainder := num % 62
		base62 = append([]byte{base62Chars[remainder]}, base62...)
		num /= 62
	}

	rand.Seed(time.Now().UnixNano())
	for len(base62) < length {
		randomChar := base62Chars[rand.Intn(len(base62Chars))]
		base62 = append([]byte{randomChar}, base62...)
	}
	return string(base62)
}
