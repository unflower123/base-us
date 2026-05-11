package batcherx

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBatcher(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		b := New()
		assert.Equal(t, 100, b.opts.size)
		assert.Equal(t, 100, b.opts.buffer)
		assert.Equal(t, 5, b.opts.worker)
		assert.Equal(t, time.Second, b.opts.interval)
	})

	t.Run("custom options", func(t *testing.T) {
		b := New(
			WithSize(200),
			WithBuffer(50),
			WithWorker(10),
			WithInterval(2*time.Second),
		)
		assert.Equal(t, 200, b.opts.size)
		assert.Equal(t, 50, b.opts.buffer)
		assert.Equal(t, 10, b.opts.worker)
		assert.Equal(t, 2*time.Second, b.opts.interval)
	})
}

func TestBatcher_Add(t *testing.T) {
	b := New(WithBuffer(2), WithWorker(2), WithBuffer(1))
	b.Do = func(ctx context.Context, val map[string][]interface{}) {

	}
	b.Start()
	defer b.Close()

	t.Run("add successfully", func(t *testing.T) {
		err := b.Add("key1", "value1")
		assert.NoError(t, err)
	})

	t.Run("channel full", func(t *testing.T) {
		// Fill the channel
		b.Add("key1", "value1")
		b.Add("key2", "value2")

		err := b.Add("key3", "value3")
		assert.Equal(t, ErrFull, err)
	})
}

func TestBatcher_merge(t *testing.T) {
	t.Run("batch by size", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		processed := make(map[string][]interface{})

		b := New(
			WithSize(2),
			WithBuffer(10),
			WithWorker(1),
			WithInterval(time.Minute), // long enough to not trigger
		)
		b.Do = func(ctx context.Context, val map[string][]interface{}) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			for k, v := range val {
				processed[k] = append(processed[k], v...)
			}
		}

		b.Start()
		defer b.Close()

		wg.Add(1) // Expect one batch of 2 items
		b.Add("key1", "value1")
		b.Add("key2", "value2")
		wg.Wait()

		assert.Equal(t, []interface{}{"value1"}, processed["key1"])
		assert.Equal(t, []interface{}{"value2"}, processed["key2"])
	})

	t.Run("batch by interval", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		processed := make(map[string][]interface{})

		b := New(
			WithSize(10), // large enough to not trigger
			WithBuffer(10),
			WithWorker(1),
			WithInterval(100*time.Millisecond),
		)
		b.Do = func(ctx context.Context, val map[string][]interface{}) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			for k, v := range val {
				processed[k] = append(processed[k], v...)
			}
		}

		b.Start()
		defer b.Close()

		wg.Add(1) // Expect one batch after interval
		b.Add("key1", "value1")
		wg.Wait()

		assert.Equal(t, []interface{}{"value1"}, processed["key1"])
	})
}

func TestBatcher_Close(t *testing.T) {
	var processed int
	var wg sync.WaitGroup

	b := New(
		WithSize(10),
		WithBuffer(10),
		WithWorker(1),
		WithInterval(time.Minute),
	)
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		defer wg.Done()
		processed += len(val["key1"])
	}

	b.Start()

	wg.Add(1)
	b.Add("key1", "value1")
	b.Add("key1", "value2")
	b.Close() // Should trigger processing of remaining items
	wg.Wait()

	assert.Equal(t, 2, processed)
}

func TestFindLightestWorker(t *testing.T) {
	b := New(WithWorker(3))
	b.Do = func(ctx context.Context, val map[string][]interface{}) {

	}
	b.Start()
	defer b.Close()

	// Initial state
	idx, ch := b.findLightestWorker()
	assert.True(t, idx >= 0 && idx < 3)
	assert.Equal(t, int32(0), atomic.LoadInt32(&ch.len))

	// Simulate adding to one channel
	atomic.AddInt32(&b.chans[0].len, 3)
	atomic.AddInt32(&b.chans[1].len, 1)
	atomic.AddInt32(&b.chans[2].len, 2)

	idx, ch = b.findLightestWorker()
	assert.Equal(t, 1, idx)
	assert.Equal(t, int32(1), atomic.LoadInt32(&ch.len))
}

func TestPanicWhenDoNil(t *testing.T) {
	b := New()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	b.Start() // Should panic because Do is nil
}
