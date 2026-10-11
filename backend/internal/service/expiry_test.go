package service_test

import (
	"testing"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
)

func TestExpireStaleHoldBookings(t *testing.T) {
	database, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	room := model.Room{
		Name:              "Executive Suite",
		Type:              "suite",
		BasePricePerNight: 2000000,
		TotalUnits:        3,
		IsActive:          true,
	}
	database.Create(&room)

	guest := model.Guest{
		Name:        "Alice",
		PhoneNumber: "089999999",
	}
	database.Create(&guest)

	// 1. Expired hold booking (expired 10 minutes ago)
	pastTime := time.Now().Add(-10 * time.Minute)
	expiredBooking := model.Booking{
		BookingCode:   "HTL-TEST-EXPIRED",
		GuestID:       guest.ID,
		RoomID:        room.ID,
		CheckInDate:   "2026-11-01",
		CheckOutDate:  "2026-11-03",
		TotalNights:   2,
		TotalRooms:    1,
		TotalPrice:    4000000,
		Status:        model.StatusHold,
		HoldExpiresAt: &pastTime,
	}
	database.Create(&expiredBooking)

	// 2. Active hold booking (expires in 1 hour)
	futureTime := time.Now().Add(1 * time.Hour)
	activeBooking := model.Booking{
		BookingCode:   "HTL-TEST-ACTIVE",
		GuestID:       guest.ID,
		RoomID:        room.ID,
		CheckInDate:   "2026-11-05",
		CheckOutDate:  "2026-11-07",
		TotalNights:   2,
		TotalRooms:    1,
		TotalPrice:    4000000,
		Status:        model.StatusHold,
		HoldExpiresAt: &futureTime,
	}
	database.Create(&activeBooking)

	// Run expiration worker logic
	expiredCount, err := service.ExpireStaleHoldBookings(database)
	if err != nil {
		t.Fatalf("failed to expire stale bookings: %v", err)
	}
	if expiredCount != 1 {
		t.Fatalf("expected 1 booking to be expired, got %d", expiredCount)
	}

	var updatedExpired model.Booking
	database.First(&updatedExpired, expiredBooking.ID)
	if updatedExpired.Status != model.StatusCancelled {
		t.Errorf("expected status CANCELLED, got %s", updatedExpired.Status)
	}

	var updatedActive model.Booking
	database.First(&updatedActive, activeBooking.ID)
	if updatedActive.Status != model.StatusHold {
		t.Errorf("expected active booking to remain HOLD, got %s", updatedActive.Status)
	}
}
