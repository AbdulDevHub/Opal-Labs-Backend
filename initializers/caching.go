package initializers

import (
	"log"
	"os"
	"strconv"

	"github.com/opalescencelabs/backend/api/caching"
)

// InitializeRedis sets up the optional Redis cache. Without REDIS_ADDR the app runs uncached.
func InitializeRedis() {
	rdAddr := os.Getenv("REDIS_ADDR")
	if rdAddr == "" {
		log.Println("REDIS_ADDR not set, running without cache")
		return
	}
	log.Println("Connecting to Redis cache...")
	rdDB := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Invalid REDIS_DB %q, using database 0", v)
		} else {
			rdDB = n
		}
	}
	if err := caching.InitializeClient(rdAddr, os.Getenv("REDIS_PASSWORD"), rdDB); err != nil {
		// Not fatal: the app still works without the cache.
		log.Printf("Error connecting to Redis, running without cache: %v", err)
		return
	}
	log.Println("Redis cache connected successfully")
}
