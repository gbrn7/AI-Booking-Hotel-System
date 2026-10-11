package db_test

import (
	"testing"

	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
)

func TestInitSQLiteAndSeed(t *testing.T) {
	database, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}

	err = db.SeedDefaultData(database)
	if err != nil {
		t.Fatalf("failed to seed default data: %v", err)
	}

	var roomCount int64
	database.Model(&model.Room{}).Count(&roomCount)
	if roomCount != 4 {
		t.Errorf("expected 4 seeded rooms, got %d", roomCount)
	}

	var userCount int64
	database.Model(&model.User{}).Count(&userCount)
	if userCount != 2 {
		t.Errorf("expected 2 seeded users, got %d", userCount)
	}

	// Verify idempotency
	err = db.SeedDefaultData(database)
	if err != nil {
		t.Fatalf("expected re-seeding to be idempotent, got error: %v", err)
	}
	database.Model(&model.Room{}).Count(&roomCount)
	if roomCount != 4 {
		t.Errorf("expected count to remain 4 after re-seed, got %d", roomCount)
	}
}
