package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes PostgreSQL or SQLite based on configuration and runs AutoMigrate.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if cfg.DBDriver == "sqlite" {
		dialector = sqlite.Open("hotel.db")
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		dialector = postgres.Open(dsn)
	}

	database, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := AutoMigrate(database); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// InitTestDB initializes an in-memory SQLite database for unit and integration testing.
func InitTestDB() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(":memory:?_pragma=foreign_keys(1)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := AutoMigrate(database); err != nil {
		return nil, err
	}

	return database, nil
}

// AutoMigrate migrates all schema tables.
func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(
		&model.Room{},
		&model.Guest{},
		&model.Booking{},
		&model.Conversation{},
		&model.Message{},
		&model.User{},
	)
}
