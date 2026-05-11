package captcha

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// RedisStore Redis存储实现
type RedisStore struct {
	client        *redis.Client
	expireSeconds int
}

// NewRedisStore 创建Redis存储实例
func NewRedisStore(client *redis.Client, expireSeconds int) *RedisStore {
	return &RedisStore{
		client:        client,
		expireSeconds: expireSeconds,
	}
}

func (s *RedisStore) Set(id string, value string) error {
	return s.client.Set(context.Background(),
		fmt.Sprintf("captcha:%s", id),
		value,
		time.Duration(s.expireSeconds)*time.Second,
	).Err()
}

func (s *RedisStore) Get(id string, clear bool) string {
	key := fmt.Sprintf("captcha:%s", id)
	val, err := s.client.Get(context.Background(), key).Result()
	if err != nil {
		return ""
	}

	if clear {
		s.client.Del(context.Background(), key)
	}

	return val
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	//return s.Get(id, clear) == answer
	storedAnswer := s.Get(id, clear)
	return strings.EqualFold(
		strings.TrimSpace(answer),
		strings.TrimSpace(storedAnswer),
	)
}
