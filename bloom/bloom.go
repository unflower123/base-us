/**
 * @Author: darry
 * @Desc:
 * @Date: 2025/4/15 11:25
 */
package bloom

import (
	"base/consts"
	"context"
	"fmt"
	"time"

	"base/rlock"

	"github.com/redis/go-redis/v9"
)

// BloomManager manages orders using a Bloom Filter and a distributed lock.
type BloomManager struct {
	bloomFilter *rlock.RedisBloomFilter
	lock        *rlock.LockForRedis
}

// NewBloomManager creates a new OrderManager.
func NewBloomManager(client *redis.Client, lockSeconds, lockRenewal time.Duration) (*BloomManager, error) {
	bf := rlock.NewRedisBloomFilter(client)
	lock, err := rlock.NewLockForRedis(context.Background(), client, lockSeconds, lockRenewal)
	if err != nil {
		return nil, err
	}
	return &BloomManager{
		bloomFilter: bf,
		lock:        lock,
	}, nil
}

// AddOrder adds an order to the Bloom Filter.
func (om *BloomManager) AddOrder(ctx context.Context, merchantId string, orderId ...string) error {
	key := fmt.Sprintf(consts.ORDER_BLOOM_ID_KEY, merchantId)
	return om.bloomFilter.Add(ctx, key, orderId)
}

// CheckOrder checks if an order exists in the Bloom Filter.
func (om *BloomManager) CheckOrder(ctx context.Context, merchantId, orderId string) (bool, error) {
	key := fmt.Sprintf(consts.ORDER_BLOOM_ID_KEY, merchantId)
	return om.bloomFilter.IsMember(ctx, key, orderId)
}

// AddMerchant adds an order to the Bloom Filter.
func (om *BloomManager) AddMerchant(ctx context.Context, merchantId ...string) error {

	return om.bloomFilter.Add(ctx, consts.MERCHANT_BLOOM_ID_KEY, merchantId)
}

// CheckMerchant checks if an order exists in the Bloom Filter.
func (om *BloomManager) CheckMerchant(ctx context.Context, merchantId string) (bool, error) {

	return om.bloomFilter.IsMember(ctx, consts.MERCHANT_BLOOM_ID_KEY, merchantId)
}

type BloomManagerV1 struct {
	bloomFilter *rlock.RedisBloomFilterV1
	key         string
}

func NewBloomManagerV1(client *redis.Client, orderType int) (*BloomManagerV1, error) {
	bf := rlock.NewRedisBloomFilterV1(client)
	key := consts.BLOOM_PAYIN_ORDER_KEY
	if orderType == 2 {
		key = consts.BLOOM_PAYOUT_ORDER_KEY
	}

	err := bf.InitBFKey(context.Background(), key)
	if err != nil {
		return nil, err
	}

	return &BloomManagerV1{
		bloomFilter: bf,
		key:         key,
	}, nil
}

// AddOrder adds an order to the Bloom Filter.
func (om *BloomManagerV1) AddOrder(ctx context.Context, orderId string) error {
	return om.bloomFilter.Add(ctx, om.key, orderId)
}

func (om *BloomManagerV1) CheckOrder(ctx context.Context, orderId string) (bool, error) {

	return om.bloomFilter.IsExists(ctx, om.key, orderId)
}
