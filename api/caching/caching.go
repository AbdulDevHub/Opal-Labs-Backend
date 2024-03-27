package caching

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitializeClient initializes the Redis client.
func InitializeClient(addr string, passwd string, db int) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	return nil
}

// Store stores the response in the cache.
func Store(c *gin.Context, request string, response interface{}) error {
	// Store the response in the cache
	err := RDB.Set(c, request, response, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Check checks if the response is already cached.
// request: the request string to be used as the key for the cache
// returns the cached response if it exists, otherwise returns an error
func Check(c *gin.Context, request string) (interface{}, error) {
	// Check if the response is already cached
	cachedResponse, err := RDB.Get(c, request).Result()
	if err == nil {
		// Return the cached response
		return cachedResponse, nil
	}
	// Response is not cached
	return nil, errors.New("Response for " + request + " not found in cache")
}

// Invalidate deletes a key from the cache
// request: the request string to be used as the key for the cache
// returns an error if the key could not be deleted
func Invalidate(c *gin.Context, request string) error {
	return RDB.Del(c, request).Err()
}
