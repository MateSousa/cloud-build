package internal

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DBHandler struct {
	db *gorm.DB
}

var dsn = "host=172.17.0.3 user=root password=root dbname=cloud-build port=5432 sslmode=disable"
var db *gorm.DB

// Initialize creates a new database connection and sets up the DBHandler.
func InitializePostgresDB() error {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return nil
}


