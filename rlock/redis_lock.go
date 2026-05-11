package rlock

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"sync"
	"time"
)

type Unlock func()

const (
	UnlockScript = `
    if redis.call('get', KEYS[1]) == ARGV[1]
    then
      return redis.call('del', KEYS[1])
     else
        return 0
     end
  `

	RenewScript = `
    if redis.call('get', KEYS[1]) == ARGV[1]
    then
      return redis.call('expire', KEYS[1], ARGV[2])
    else
      return 0
    end
  `
)

type LockForRedis struct {
	store    *redis.Client
	seconds  time.Duration
	renewal  time.Duration
	unlockCh chan struct{}
	once     sync.Once
}

func NewLockForRedis(ctx context.Context, store *redis.Client, seconds, renewal time.Duration) (*LockForRedis, error) {
	_, err := store.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("redis ping err: %v", err.Error()))
	}

	if seconds < renewal {
		return nil, errors.New("renewal must be less than seconds")
	}
	return &LockForRedis{
		store:    store,
		seconds:  time.Second * seconds,
		renewal:  time.Second * renewal,
		unlockCh: make(chan struct{}),
	}, nil
}

func (lock *LockForRedis) Lock(ctx context.Context, key string, options ...LockOption) (Unlock, error) {
	config := &LockConfig{
		RetryCount: 0,
		RetryDelay: 100 * time.Millisecond,
		Timeout:    0,
		maxRenewal: 0,
	}

	for _, opt := range options {
		opt(config)
	}

	var retryCount int
	startTime := time.Now()
	value := uuid.New().String()
	for {
		lockSuccess, err := lock.store.SetNX(ctx, key, value, lock.seconds).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to acquire lock: %w", err)
		}

		if lockSuccess {
			go lock.watchDog(ctx, key, value, config.maxRenewal)
			return func() {
				script := redis.NewScript(UnlockScript)
				resp := script.Run(ctx, lock.store, []string{key}, value)
				result, err := resp.Result()
				if err != nil {
					log.Printf("unlock failed: %v", err)
					return
				}

				if result == int64(0) {
					log.Println("unlock failed: key does not match")
					return
				}

				lock.once.Do(func() {
					close(lock.unlockCh)
				})
				return
			}, nil
		}

		if config.RetryCount == 0 && config.Timeout == 0 {
			return nil, errors.New("lock is already held")
		}

		if config.RetryCount > 0 && retryCount >= config.RetryCount {
			return nil, fmt.Errorf("failed to acquire lock after %d retries", config.RetryCount)
		}

		if config.Timeout > 0 && time.Since(startTime) >= config.Timeout {
			return nil, fmt.Errorf("failed to acquire lock within %v", config.Timeout)
		}

		time.Sleep(config.RetryDelay)
		retryCount++
	}
}

func (lock *LockForRedis) watchDog(ctx context.Context, key string, value string, maxRenewal int) {
	expTicker := time.NewTicker(lock.renewal)
	defer expTicker.Stop()

	script := redis.NewScript(RenewScript)
	renewalCount := 0

	for {
		select {
		case <-expTicker.C:
			if renewalCount >= maxRenewal && maxRenewal != 0 {
				log.Println("reached max renewal count, stopping watchdog")
				return
			}

			resp := script.Run(ctx, lock.store, []string{key}, value, int(lock.seconds.Seconds()))
			result, err := resp.Result()
			if err != nil {
				log.Printf("renewal failed: %v", err)
				return
			}

			if result == int64(0) {
				log.Println("renewal failed: key does not match")
				return
			}

			renewalCount++
			log.Printf("lock renewed, count: %d/%d", renewalCount, maxRenewal)

		case <-lock.unlockCh:
			log.Println("watchdog stopped by unlock")
			return
		}
	}
}

type LockConfig struct {
	RetryCount int
	RetryDelay time.Duration
	Timeout    time.Duration
	maxRenewal int
}

type LockOption func(*LockConfig)

func WithRetryCount(retryCount int) LockOption {
	return func(c *LockConfig) {
		c.RetryCount = retryCount
	}
}

func WithRetryDelay(retryDelay time.Duration) LockOption {
	return func(c *LockConfig) {
		c.RetryDelay = retryDelay
	}
}

func WithTimeout(timeout time.Duration) LockOption {
	return func(c *LockConfig) {
		c.Timeout = timeout
	}
}

func WithMaxRenewal(maxRenewal int) LockOption {
	return func(c *LockConfig) {
		c.maxRenewal = maxRenewal
	}
}
