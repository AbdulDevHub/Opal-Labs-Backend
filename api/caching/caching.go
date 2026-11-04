package caching

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RDB is nil when no Redis is configured; every helper below is then a no-op.
var RDB *redis.Client

// InitializeClient connects to Redis and pings it, leaving RDB nil on failure.
func InitializeClient(addr string, passwd string, db int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return err
	}
	RDB = client
	return nil
}

// Store stores the response in the cache.
func Store(c *gin.Context, request string, response interface{}) error {
	if RDB == nil {
		return nil
	}
	return RDB.Set(c, request, response, 0).Err()
}

// Check checks if the response is already cached.
// request: the request string to be used as the key for the cache
// returns the cached response if it exists, otherwise returns an error
func Check(c *gin.Context, request string) (interface{}, error) {
	if RDB != nil {
		if cachedResponse, err := RDB.Get(c, request).Result(); err == nil {
			return cachedResponse, nil
		}
	}
	return nil, errors.New("Response for " + request + " not found in cache")
}

// Invalidate deletes a key from the cache
// request: the request string to be used as the key for the cache
// returns an error if the key could not be deleted
func Invalidate(c *gin.Context, request string) error {
	if RDB == nil {
		return nil
	}
	return RDB.Del(c, request).Err()
}
