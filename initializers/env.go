package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvVariables loads environment variables from a .env file if one exists.
// Hosts like Render inject variables directly, so a missing file is not an error.
func LoadEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	// Without a secret, session tokens would be signed with an empty key
	if os.Getenv("SECRET") == "" {
		log.Fatal("SECRET environment variable is not set")
	}
}
