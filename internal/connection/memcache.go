package connection

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"Homework-1/internal/config"
	"Homework-1/pkg/errlst"
)

const numShards = 32

// Cache is
var _ Cache = (*InMemoryCache)(nil)

// CachedValue is
type CachedValue struct {
	data           []byte
	expirationTime time.Time
}

// InMemoryCache is
type cacheShard struct {
	cache   map[string]CachedValue
	mxCache sync.RWMutex
}

// InMemoryCache is
type InMemoryCache struct {
	shards [numShards]*cacheShard
}

// NewInMemoryCache is
func NewInMemoryCache(ctx context.Context, cfgs config.InMemoryCache) *InMemoryCache {
	cache := &InMemoryCache{}
	for i := 0; i < numShards; i++ {
		cache.shards[i] = &cacheShard{
			cache:   make(map[string]CachedValue),
			mxCache: sync.RWMutex{},
		}
	}

	go cache.evictExpiredEntries(ctx, cfgs.CleanTime)

	return cache
}

func (i *InMemoryCache) evictExpiredEntries(ctx context.Context, cleanTime float64) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			for _, shard := range i.shards {
				shard.mxCache.Lock()
				for key, v := range shard.cache {
					if now.Sub(v.expirationTime).Hours() > cleanTime {
						delete(shard.cache, key)
					}
				}
				shard.mxCache.Unlock()
			}
		default:
			time.Sleep(10 * time.Minute)
		}
	}
}

func (i *InMemoryCache) getShard(key string) *cacheShard {
	hasher := fnv.New32()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return nil
	}
	return i.shards[hasher.Sum32()%numShards]
}

// Get is
func (i *InMemoryCache) Get(_ context.Context, key string) (string, error) {
	shard := i.getShard(key)
	if shard == nil {
		return "", errlst.ErrInMemoryCacheNil
	}
	shard.mxCache.RLock()
	value, ok := shard.cache[key]
	shard.mxCache.RUnlock()

	if !ok {
		return "", errlst.ErrInMemoryCacheNil
	}

	if value.expirationTime.Before(time.Now()) {
		shard.mxCache.Lock()
		delete(shard.cache, key)
		shard.mxCache.Unlock()
		return "", errlst.ErrInMemoryCacheNil
	}

	return string(value.data), nil
}

// Set is
func (i *InMemoryCache) Set(_ context.Context, key string, value string, expiration time.Duration) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}
	shard := i.getShard(key)
	if shard == nil {
		return errlst.ErrInMemoryCacheNil
	}
	shard.mxCache.Lock()
	shard.cache[key] = CachedValue{
		data:           []byte(value),
		expirationTime: time.Now().Add(expiration),
	}
	shard.mxCache.Unlock()
	return nil
}

// Del is
func (i *InMemoryCache) Del(_ context.Context, key string) error {
	shard := i.getShard(key)
	if shard == nil {
		return errlst.ErrInMemoryCacheNil
	}
	shard.mxCache.Lock()
	delete(shard.cache, key)
	shard.mxCache.Unlock()
	return nil
}

// Close is
func (i *InMemoryCache) Close() error {
	for _, shard := range i.shards {
		shard.mxCache.Lock()
		shard.cache = make(map[string]CachedValue)
		shard.mxCache.Unlock()
	}
	return nil
}
