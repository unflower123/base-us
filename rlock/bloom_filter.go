/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/15 09:42
 */
package rlock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
)

//
//var (
//	capacity  = 4194304000
//	errorRate = 0.01
//)
//
//// RedisBloomFilter is a Bloom Filter using Redis Bitmap.
//type RedisBloomFilter struct {
//	client *redis.Client
//	//key       string
//	capacity  int     // capacity max number 4194304000
//	errorRate float64 // Desired false positive rate 0.01%
//	size      int     // Bitmap size in bits
//	hashCount int     // Number of hash functions
//}
//
//// NewRedisBloomFilter creates a new Bloom Filter with the given parameters.
//// - ctx: Context for managing request lifecycle (optional)
//// - client: Redis client
//func NewRedisBloomFilter(client *redis.Client) *RedisBloomFilter {
//	bf := &RedisBloomFilter{
//		client:    client,
//		capacity:  capacity,
//		errorRate: errorRate,
//	}
//	bf.size = bf.optimalSize()
//	bf.hashCount = bf.optimalHashCount()
//	return bf
//}
//
//// optimalSize calculates the optimal bitmap size (bits).
//// m = - (n * ln(p)) / (ln(2)^2)
//func (bf *RedisBloomFilter) optimalSize() int {
//	m := -float64(bf.capacity) * math.Log(bf.errorRate) / (math.Ln2 * math.Ln2)
//	return int(math.Ceil(m))
//}
//
//// optimalHashCount calculates the optimal number of hash functions.
//// k = (m/n) * ln(2)
//func (bf *RedisBloomFilter) optimalHashCount() int {
//	k := float64(bf.size) / float64(bf.capacity) * math.Ln2
//	return int(math.Max(1, math.Ceil(k)))
//}
//
//// hash generates a hash for an item using Murmur3 with a seed.
//func (bf *RedisBloomFilter) hash(item string, seed int) int {
//	// Create a Murmur3 hash with the given seed
//	hasher := murmur3.New32WithSeed(uint32(seed))
//	// Write the item to the hasher
//	hasher.Write([]byte(item))
//	// Sum the hash and convert to int
//	hash := hasher.Sum32()
//	// Return the hash modulo the size of the bitmap
//	return int(hash % uint32(bf.size))
//}
//
//// Add adds an item to the Bloom Filter.
//func (bf *RedisBloomFilter) Add(ctx context.Context, key, value string) error {
//	return bf.addItems(ctx, key, []string{value})
//}
//
//// AddMulti adds multiple items to the Bloom Filter.
//func (bf *RedisBloomFilter) AddMulti(ctx context.Context, key string, values []string) error {
//	return bf.addItems(ctx, key, values)
//}
//
//func (bf *RedisBloomFilter) addItems(ctx context.Context, key string, values []string) error {
//	pipe := bf.client.Pipeline()
//	for _, item := range values {
//		for i := 0; i < bf.hashCount; i++ {
//			bitIndex := bf.hash(item, i)
//			fmt.Printf("index = = = %v,bitIndex = %v \n", i, bitIndex)
//			pipe.SetBit(ctx, key, int64(bitIndex), 1)
//		}
//	}
//	_, err := pipe.Exec(ctx)
//	return err
//}
//
//// Check checks if an item might exist in the Bloom Filter.
//// Returns true if possibly present, false if definitely not.
//func (bf *RedisBloomFilter) Check(ctx context.Context, key string, values string) (bool, error) {
//	for i := 0; i < bf.hashCount; i++ {
//		bitIndex := bf.hash(values, i)
//		bit, err := bf.client.GetBit(ctx, key, int64(bitIndex)).Result()
//		if err != nil {
//			return false, err
//		}
//		if bit == 0 {
//			return false, nil
//		}
//	}
//	return true, nil
//}
//
//// CheckMulti checks multiple items in the Bloom Filter.
//func (bf *RedisBloomFilter) CheckMulti(ctx context.Context, key string, values []string) ([]bool, error) {
//	results := make([]bool, len(values))
//	for i, item := range values {
//		exists, err := bf.Check(ctx, key, item)
//		if err != nil {
//			return nil, err
//		}
//		results[i] = exists
//	}
//	return results, nil
//}
//
//// Clear removes the Bloom Filter from Redis.
//func (bf *RedisBloomFilter) Clear(ctx context.Context, key string) error {
//	return bf.client.Del(ctx, key).Err()
//}

// RedisSet is a wrapper around Redis Set operations for basic usage.
type RedisBloomFilter struct {
	client *redis.Client
}

// NewRedisSet creates a new RedisSet instance for a given key.
func NewRedisBloomFilter(client *redis.Client) *RedisBloomFilter {
	return &RedisBloomFilter{
		client: client,
	}
}

// Add adds one or more members to the set.
func (s *RedisBloomFilter) Add(ctx context.Context, key string, members ...any) error {
	return s.client.SAdd(ctx, key, members...).Err()
}

// IsMember checks if a member exists in the set.
func (s *RedisBloomFilter) IsMember(ctx context.Context, key string, member any) (bool, error) {
	return s.client.SIsMember(ctx, key, member).Result()
}

// Members returns all members of the set.
func (s *RedisBloomFilter) Members(ctx context.Context, key string) ([]string, error) {
	return s.client.SMembers(ctx, key).Result()
}

// Remove removes one or more members from the set.
func (s *RedisBloomFilter) Remove(ctx context.Context, key string, members ...any) error {
	return s.client.SRem(ctx, key, members...).Err()
}

// Clear removes all members from the set (empties the set).
func (s *RedisBloomFilter) Clear(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// RemoveMember removes a specific member from the set.
func (s *RedisBloomFilter) RemoveMember(ctx context.Context, key string, member any) error {
	return s.client.SRem(ctx, key, member).Err()
}

// Len returns the cardinality (number of elements) of the set.
func (s *RedisBloomFilter) Len(ctx context.Context, key string) (int64, error) {
	return s.client.SCard(ctx, key).Result()
}

type RedisBloomFilterV1 struct {
	client    *redis.Client
	errorRate float64
	capactity int64
}

// NewRedisSet creates a new RedisSet instance for a given key.
func NewRedisBloomFilterV1(client *redis.Client) *RedisBloomFilterV1 {
	return &RedisBloomFilterV1{
		client:    client,
		errorRate: 0.01,
		//capactity: 10000000000,
		capactity: 1073741824,
	}
}

func (s *RedisBloomFilterV1) InitBFKey(ctx context.Context, key string) error {
	err := s.client.BFReserve(ctx, key, s.errorRate, s.capactity).Err()
	if err != nil && strings.Contains(err.Error(), "ERR item exists") {
		return nil
	}

	return err
}

func (s *RedisBloomFilterV1) Add(ctx context.Context, key string, members any) error {
	return s.client.BFAdd(ctx, key, members).Err()
}

func (s *RedisBloomFilterV1) Adds(ctx context.Context, key string, members ...any) error {
	return s.client.BFMAdd(ctx, key, members).Err()
}

func (s *RedisBloomFilterV1) IsExists(ctx context.Context, key string, member any) (bool, error) {
	return s.client.BFExists(ctx, key, member).Result()
}

func (s *RedisBloomFilterV1) IsMExists(ctx context.Context, key string, member ...any) ([]bool, error) {
	return s.client.BFMExists(ctx, key, member...).Result()
}
