package bloom

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"base/consts"
)

// TestBloomManagerAddOrder tests the AddOrder method of BloomManager.
func TestBloomManagerAddOrder_a(t *testing.T) {
	// Start a miniredis server

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test AddOrder
	ctx := context.Background()
	merchantID := "merchant123"
	orderID := "order456-789"

	// AddOrder
	err = bloomManager.AddOrder(ctx, merchantID, orderID)
	assert.NoError(t, err)

	// CheckOrder to verify the order was added
	exists, err := bloomManager.CheckOrder(ctx, merchantID, orderID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check a non-existent order
	nonExistentOrderID := "order789"
	exists, err = bloomManager.CheckOrder(ctx, merchantID, nonExistentOrderID)
	assert.NoError(t, err)
	assert.False(t, exists)

	key := fmt.Sprintf(consts.ORDER_BLOOM_ID_KEY, merchantID)
	bloomManager.bloomFilter.Clear(ctx, key)
}

// TestBloomManagerAddOrder tests the AddOrder method of BloomManager.
func TestBloomManagerAddOrder(t *testing.T) {
	// Start a miniredis server
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test AddOrder
	ctx := context.Background()
	merchantID := "merchant123"
	orderID := "order456"

	// AddOrder
	err = bloomManager.AddOrder(ctx, merchantID, orderID)
	assert.NoError(t, err)

	// CheckOrder to verify the order was added
	exists, err := bloomManager.CheckOrder(ctx, merchantID, orderID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check a non-existent order
	nonExistentOrderID := "order789"
	exists, err = bloomManager.CheckOrder(ctx, merchantID, nonExistentOrderID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestBloomManagerConcurrency tests the concurrency of BloomManager methods.
func TestBloomManagerConcurrency(t *testing.T) {
	// Start a miniredis server
	//srv, err := miniredis.Run()
	//require.NoError(t, err)
	//defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		//Addr: srv.Addr(),
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test concurrent AddOrder and CheckOrder
	ctx := context.Background()
	merchantID := "merchant123"
	orderID := "order456"

	// AddOrder concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := bloomManager.AddOrder(ctx, merchantID, orderID)
			assert.NoError(t, err)
			t.Logf("= = = = = index : %v,err : %v\n", i, err)
		}(i)
	}
	wg.Wait()

	//CheckOrder concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			exists, err := bloomManager.CheckOrder(ctx, merchantID, orderID)
			assert.NoError(t, err)
			assert.True(t, exists)
			t.Logf("= = = = = check index : %v,err : %v\n", i, err)
		}(i)
	}
	wg.Wait()

	// Test Clear
	key := fmt.Sprintf(consts.ORDER_BLOOM_ID_KEY, merchantID)
	err = bloomManager.bloomFilter.Clear(ctx, key)
	assert.NoError(t, err)

	var res int64
	// Verify that the key is deleted
	res, err = client.Exists(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}

// TestBloomManagerTimeout tests the behavior when Redis operations time out.
func TestBloomManagerTimeout(t *testing.T) {
	// Start a miniredis server
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr:        srv.Addr(),
		ReadTimeout: 1 * time.Nanosecond, // Set a very short read timeout to simulate timeout
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test AddOrder with timeout
	ctx := context.Background()
	merchantID := "merchant123"
	orderID := "order456"
	err = bloomManager.AddOrder(ctx, merchantID, orderID)
	assert.Error(t, err)

	// Test CheckOrder with timeout
	exists, err := bloomManager.CheckOrder(ctx, merchantID, orderID)
	assert.Error(t, err)
	assert.False(t, exists)
}

// TestBloomManagerLargeData tests the BloomManager with a large amount of data.
func TestBloomManagerLargeData(t *testing.T) {
	// Start a miniredis server
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test AddOrder and CheckOrder with a large number of items
	ctx := context.Background()
	merchantID := "merchant123"
	largeItems := generateLargeItems(10)

	// AddOrder concurrently
	var wg sync.WaitGroup
	for _, item := range largeItems {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			err := bloomManager.AddOrder(ctx, merchantID, item)
			assert.NoError(t, err)
		}(item)
	}
	wg.Wait()
	largeItems = generateLargeItems(10)
	// CheckOrder concurrently
	for _, item := range largeItems {
		wg.Add(1)
		go func(item string) {
			defer wg.Done()
			exists, err := bloomManager.CheckOrder(ctx, merchantID, item)
			assert.NoError(t, err)
			assert.True(t, exists)
		}(item)
	}
	wg.Wait()

	// Test Clear
	key := fmt.Sprintf(consts.ORDER_BLOOM_ID_KEY, merchantID)
	err = bloomManager.bloomFilter.Clear(ctx, key)
	assert.NoError(t, err)

	var res int64
	// Verify that the key is deleted
	res, err = client.Exists(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}

// generateLargeItems generates a slice of large items.
func generateLargeItems(n int) []string {
	items := make([]string, n)
	for i := 0; i < n; i++ {
		items[i] = fmt.Sprintf("order%d", i)
	}
	return items
}

// TestBloomManagerAddMerchant tests the AddMerchant and CheckMerchant methods of BloomManager.
func TestBloomManagerAddMerchant(t *testing.T) {
	// Start a miniredis server
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test AddMerchant
	ctx := context.Background()
	merchantID := "merchant123"

	// AddMerchant
	err = bloomManager.AddMerchant(ctx, merchantID)
	assert.NoError(t, err)

	// CheckMerchant to verify the merchant was added
	exists, err := bloomManager.CheckMerchant(ctx, merchantID)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check a non-existent merchant
	nonExistentMerchantID := "merchant456"
	exists, err = bloomManager.CheckMerchant(ctx, nonExistentMerchantID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

// TestBloomManagerConcurrencyMerchant tests the concurrency of AddMerchant and CheckMerchant methods.
func TestBloomManagerConcurrencyMerchant(t *testing.T) {
	// Start a miniredis server
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	// Create a Redis client connected to the miniredis server
	client := redis.NewClient(&redis.Options{
		Addr: srv.Addr(),
	})

	// Create a new BloomManager
	lockSeconds := time.Duration(10)
	lockRenewal := time.Duration(5)
	bloomManager, err := NewBloomManager(client, lockSeconds, lockRenewal)
	require.NoError(t, err)

	// Test concurrent AddMerchant and CheckMerchant
	ctx := context.Background()
	merchantIDs := []string{"merchant1", "merchant2", "merchant3"}

	// AddMerchant concurrently
	var wg sync.WaitGroup
	for _, merchantID := range merchantIDs {
		wg.Add(1)
		go func(merchantID string) {
			defer wg.Done()
			err := bloomManager.AddMerchant(ctx, merchantID)
			assert.NoError(t, err)
		}(merchantID)
	}
	wg.Wait()

	// CheckMerchant concurrently
	for _, merchantID := range merchantIDs {
		wg.Add(1)
		go func(merchantID string) {
			defer wg.Done()
			exists, err := bloomManager.CheckMerchant(ctx, merchantID)
			assert.NoError(t, err)
			assert.True(t, exists)
		}(merchantID)
	}
	wg.Wait()

	// Test Clear
	key := consts.MERCHANT_BLOOM_ID_KEY
	err = bloomManager.bloomFilter.Clear(ctx, key)
	assert.NoError(t, err)

	var res int64
	// Verify that the key is deleted
	res, err = client.Exists(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}
