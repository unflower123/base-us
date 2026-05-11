package batcherx

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var ErrFull = errors.New("channel is full")

type Option interface {
	apply(*options)
}

type options struct {
	size     int
	buffer   int
	worker   int
	interval time.Duration
}

func (o *options) check() {
	if o.size <= 0 {
		o.size = 100
	}
	if o.buffer <= 0 {
		o.buffer = 100
	}
	if o.worker <= 0 {
		o.worker = 5
	}
	if o.interval <= 0 {
		o.interval = time.Second
	}
}

type funcOption struct {
	f func(*options)
}

func (fo *funcOption) apply(o *options) {
	fo.f(o)
}

func newOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithSize(s int) Option {
	return newOption(func(o *options) {
		o.size = s
	})
}

func WithBuffer(b int) Option {
	return newOption(func(o *options) {
		o.buffer = b
	})
}

func WithWorker(w int) Option {
	return newOption(func(o *options) {
		o.worker = w
	})
}

func WithInterval(i time.Duration) Option {
	return newOption(func(o *options) {
		o.interval = i
	})
}

type msg struct {
	key string
	val interface{}
}

type channelWithLen struct {
	ch  chan *msg
	len int32
}

type Batcher struct {
	opts   options
	Do     func(ctx context.Context, val map[string][]interface{})
	chans  []*channelWithLen
	wait   sync.WaitGroup
	closed int32
}

func New(opts ...Option) *Batcher {
	b := &Batcher{}
	for _, opt := range opts {
		opt.apply(&b.opts)
	}
	b.opts.check()

	b.chans = make([]*channelWithLen, b.opts.worker)
	for i := 0; i < b.opts.worker; i++ {
		b.chans[i] = &channelWithLen{
			ch:  make(chan *msg, b.opts.buffer),
			len: 0,
		}
	}
	return b
}

func (b *Batcher) Start() {
	if b.Do == nil {
		log.Fatal("Batcher: Do func is nil")
	}
	b.wait.Add(len(b.chans))
	for i, ch := range b.chans {
		go b.merge(i, ch.ch)
	}
}

func (b *Batcher) findLightestWorker() (int, *channelWithLen) {
	minIdx := 0
	minLen := atomic.LoadInt32(&b.chans[0].len)

	for i := 1; i < len(b.chans); i++ {
		currLen := atomic.LoadInt32(&b.chans[i].len)
		if currLen < minLen {
			minIdx = i
			minLen = currLen
		}
	}

	return minIdx, b.chans[minIdx]
}

func (b *Batcher) Add(key string, val interface{}) error {
	if atomic.LoadInt32(&b.closed) == 1 {
		return errors.New("batcher is closed")
	}

	_, lightest := b.findLightestWorker()

	if atomic.LoadInt32(&lightest.len) >= int32(b.opts.buffer) {
		return ErrFull
	}

	atomic.AddInt32(&lightest.len, 1)
	select {
	case lightest.ch <- &msg{key: key, val: val}:
		return nil
	default:
		atomic.AddInt32(&lightest.len, -1)
		return ErrFull
	}
}

func (b *Batcher) merge(idx int, ch <-chan *msg) {
	defer b.wait.Done()

	var (
		count  int
		vals   = make(map[string][]interface{}, b.opts.size)
		ticker = time.NewTicker(b.opts.interval)
	)

	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				if len(vals) > 0 {
					b.Do(context.Background(), vals)
				}
				return
			}
			atomic.AddInt32(&b.chans[idx].len, -1)
			count++
			vals[msg.key] = append(vals[msg.key], msg.val)
			if count >= b.opts.size {
				b.Do(context.Background(), vals)
				vals = make(map[string][]interface{}, b.opts.size)
				count = 0
			}
		case <-ticker.C:
			if len(vals) > 0 {
				b.Do(context.Background(), vals)
				vals = make(map[string][]interface{}, b.opts.size)
				count = 0
			}
		}
	}
}

func (b *Batcher) Close() {
	if !atomic.CompareAndSwapInt32(&b.closed, 0, 1) {
		return
	}
	for _, ch := range b.chans {
		close(ch.ch)
	}
	b.wait.Wait()
}
