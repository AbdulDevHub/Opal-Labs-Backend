package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/opalescencelabs/backend/database"
)

// ConnectToDB connects to the database and migrates the tables.
// DATABASE_URL (e.g. a Neon connection string) takes priority;
// otherwise the DSN is built from the individual DB_* variables.
func ConnectToDB() {
	log.Println("Connecting to DB...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_SSL_MODE"))
	}

	if err := database.Connect(dsn); err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}
	log.Println("Database connected successfully")

	if err := database.Migrate(); err != nil {
		panic("Table migration failed: " + err.Error())
	}
	log.Println("Tables migrated successfully")
}
