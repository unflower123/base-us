package rlock

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
	"testing"
	"time"
)

func newMockRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func TestLockSuccess(t *testing.T) {
	client := newMockRedisClient()
	defer client.Close()
	ctx := context.Background()
	lock, err := NewLockForRedis(ctx, client, 10*time.Second, 5*time.Second)
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	unlock, err := lock.Lock(ctx, "test-key")
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	defer unlock()

	val, err := client.Get(ctx, "test-key").Result()
	if err != nil {
		t.Fatalf("failed to get lock value: %v", err)
	}
	if val != "test-value" {
		t.Fatalf("expected lock value to be 'test-value', got '%s'", val)
	}
	fmt.Println(val)
}

func TestLockRetry(t *testing.T) {
	client := newMockRedisClient()
	ctx := context.Background()
	defer client.Close()

	lock, err := NewLockForRedis(ctx, client, 10, 5)
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	unlock1, err := lock.Lock(ctx, "test-key")
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	defer unlock1()

	_, err = lock.Lock(ctx, "test-key", WithRetryCount(3), WithRetryDelay(100*time.Millisecond))
	if err == nil {
		t.Fatal("expected lock to be held, but it was acquired")
	}
	if !strings.Contains(err.Error(), "failed to acquire lock after 3 retries") {
		t.Fatalf("expected error 'failed to acquire lock after 3 retries', got '%v'", err)
	}
}

func TestLockRenewal(t *testing.T) {
	client := newMockRedisClient()
	defer client.Close()
	ctx := context.Background()
	lock, err := NewLockForRedis(ctx, client, 2, 1)
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	unlock, err := lock.Lock(ctx, "test-key")
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	defer unlock()

	time.Sleep(3 * time.Second)

	ttl, err := client.TTL(ctx, "test-key").Result()
	if err != nil {
		t.Fatalf("failed to get lock TTL: %v", err)
	}
	if ttl <= 0 {
		t.Fatalf("expected lock to be renewed, but TTL is %v", ttl)
	}
}

func TestLockRelease(t *testing.T) {
	client := newMockRedisClient()
	defer client.Close()
	ctx := context.Background()
	lock, err := NewLockForRedis(ctx, client, 10, 5)
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	unlock, err := lock.Lock(ctx, "test-key")
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}

	unlock()

	val, err := client.Get(ctx, "test-key").Result()
	if !errors.Is(err, redis.Nil) {
		t.Fatalf("expected lock to be released, but it still exists: %v", val)
	}
}

func TestLockMaxRenewal(t *testing.T) {
	client := newMockRedisClient()
	defer client.Close()

	ctx := context.Background()

	lock, err := NewLockForRedis(ctx, client, 2, 1)
	if err != nil {
		t.Fatalf("failed to create lock: %v", err)
	}

	unlock, err := lock.Lock(ctx, "test-key", WithMaxRenewal(2))
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}
	defer unlock()

	time.Sleep(3 * time.Second)

	ttl, err := client.TTL(ctx, "test-key").Result()
	if err != nil {
		t.Fatalf("failed to get lock TTL: %v", err)
	}
	if ttl > 0 {
		t.Fatalf("expected lock to expire, but TTL is %v", ttl)
	}
}

func TestLockRedisConnectionFailure(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "invalid-address:6379",
	})
	ctx := context.Background()
	_, err := NewLockForRedis(ctx, client, 10*time.Second, 5*time.Second)
	if err == nil {
		t.Fatal("expected error due to Redis connection failure, but got nil")
	}
}
