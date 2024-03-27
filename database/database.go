package database

import (
	"github.com/opalescencelabs/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to the database
// dsn is the Data Source Name,
// e.g. "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable".
// Returns error if connection fails, nil otherwise.
func Connect(dsn string) error {
	var err error
	// TODO: fix foreign-key-constraint error when adding new Users without pages
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return err
	}
  
	return nil
}

// Migrate the database
// AutoMigrate the Element, User, and Page models.
// Returns error if migration fails, nil otherwise.
func Migrate() error {
	var err error
	err = DB.AutoMigrate(&models.Element{}, &models.User{}, &models.Page{})
	if err != nil {
		return err
	}
	return nil
}
