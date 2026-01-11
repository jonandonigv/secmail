package database

import (
	"secmail/internal/email"
	"secmail/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the PostgreSQL database connection and auto-migrates schemas.
func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &email.Message{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
