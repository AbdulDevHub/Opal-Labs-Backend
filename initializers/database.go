package initializers

import (
	"log"

	"fmt"
	"os"
	"github.com/opalescencelabs/backend/database"
)

// ConnectToDB connects to the database and migrates the tables.
func ConnectToDB() {
	log.Println("Connecting to DB...")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	if err := database.Connect(dsn); err != nil {
		panic("Failed to connect to DB")
	}
	log.Println("Database connected successfully")

	if err := database.Migrate(); err != nil {
		panic("Table migration failed")
	}
	log.Println("Tables migrated successfully")
}
