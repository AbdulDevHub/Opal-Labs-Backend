package initializers

import (
	"log"
	"os"
	"strconv"

	"github.com/opalescencelabs/backend/api/caching"
)

// InitializeRedis initializes the Redis client.
func InitializeRedis() {
	log.Println("Connecting to Redis cache...")
	rdAddr := os.Getenv("REDIS_ADDR")
	rdPassword := os.Getenv("REDIS_PASSWORD")
	rdDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB to int: %v, using default database 0", err)
		rdDB = 0
	}
	if err := caching.InitializeClient(rdAddr, rdPassword, rdDB); err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
		// Not panic because the app can still run without cache
	} else {
		log.Println("Redis Cache connected successfully")
	}
}
