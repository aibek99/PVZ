package cache

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"Homework-1/internal/connection"
	"Homework-1/internal/model/abstract"
	"Homework-1/pkg/tracing"
)

// Store is
var _ Store = (*ClientRedisRepo)(nil)

// Store is
type Store interface {
	Get(ctx context.Context, argument abstract.CacheArgument) ([]byte, error)
	Set(ctx context.Context, argument abstract.CacheArgument, value []byte, duration time.Duration) error
	Del(ctx context.Context, argument abstract.CacheArgument) error
}

// ClientRedisRepo is
type ClientRedisRepo struct {
	rdb connection.Cache
}

// NewClientRDRepository is
func NewClientRDRepository(rdb connection.Cache) *ClientRedisRepo {
	return &ClientRedisRepo{rdb: rdb}
}

func (c *ClientRedisRepo) getCacheKey(objectType string, id string) string {
	return strings.Join([]string{
		objectType,
		id,
	}, ":")
}

// Get is
func (c *ClientRedisRepo) Get(ctx context.Context, argument abstract.CacheArgument) ([]byte, error) {
	tracer := otel.Tracer("[ClientRedisRepo]")
	ctx, span := tracer.Start(ctx, "[Get]")
	defer span.End()

	key := argument.ToCacheStorage()
	cacheKey := c.getCacheKey(key.ObjectType, key.ID)
	valueString, err := c.rdb.Get(ctx, cacheKey)
	if err != nil {
		tracing.ErrorTracer(span, err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "Successfully get data from redis repo")
	return []byte(valueString), nil
}

// Set is
func (c *ClientRedisRepo) Set(ctx context.Context, argument abstract.CacheArgument, value []byte, duration time.Duration) error {
	tracer := otel.Tracer("[ClientRedisRepo]")
	ctx, span := tracer.Start(ctx, "[Set]")
	defer span.End()

	key := argument.ToCacheStorage()
	cacheKey := c.getCacheKey(key.ObjectType, key.ID)

	err := c.rdb.Set(ctx, cacheKey, string(value), duration)
	if err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "Successfully set data from redis repo")
	return nil
}

// Del is
func (c *ClientRedisRepo) Del(ctx context.Context, argument abstract.CacheArgument) error {
	tracer := otel.Tracer("[ClientRedisRepo]")
	ctx, span := tracer.Start(ctx, "[Del]")
	defer span.End()

	key := argument.ToCacheStorage()
	cacheKey := c.getCacheKey(key.ObjectType, key.ID)

	err := c.rdb.Del(ctx, cacheKey)
	if err != nil {
		tracing.ErrorTracer(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "Successfully deleted data from redis repo")
	return nil
}
