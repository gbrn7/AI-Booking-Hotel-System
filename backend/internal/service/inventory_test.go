package service_test

import (
	"sync"
	"testing"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
)

func TestCheckAvailabilityAndBooking(t *testing.T) {
	database, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	room := model.Room{
		Name:              "Sunset Deluxe",
		Type:              "deluxe",
		BasePricePerNight: 1000000,
		CapacityGuests:    2,
		TotalUnits:        2,
		IsActive:          true,
	}
	database.Create(&room)

	checkIn := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	checkOut := time.Now().Add(72 * time.Hour).Format("2006-01-02")

	// 1. Initial check - should have 2 units available
	avails, err := service.CheckAvailability(database, checkIn, checkOut, "deluxe")
	if err != nil {
		t.Fatalf("unexpected error checking availability: %v", err)
	}
	if len(avails) != 1 || avails[0].AvailableUnits != 2 {
		t.Fatalf("expected 2 available units, got %+v", avails)
	}

	// 2. Hold 1 room
	bookingReq := service.HoldBookingRequest{
		GuestName:       "John Doe",
		GuestPhone:      "08123456789",
		RoomID:          room.ID,
		CheckInDate:     checkIn,
		CheckOutDate:    checkOut,
		TotalRooms:      1,
		HoldDurationMin: 120,
	}
	booking, err := service.HoldBooking(database, bookingReq)
	if err != nil {
		t.Fatalf("failed to hold booking: %v", err)
	}
	if booking.Status != model.StatusHold {
		t.Errorf("expected status HOLD, got %s", booking.Status)
	}
	if booking.TotalNights != 2 {
		t.Errorf("expected 2 nights, got %d", booking.TotalNights)
	}
	if booking.TotalPrice != 2000000 {
		t.Errorf("expected total price 2,000,000, got %f", booking.TotalPrice)
	}

	// 3. Re-check availability - should now have 1 unit left
	avails, err = service.CheckAvailability(database, checkIn, checkOut, "deluxe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if avails[0].AvailableUnits != 1 {
		t.Errorf("expected 1 remaining unit, got %d", avails[0].AvailableUnits)
	}
}

func TestHoldBookingConcurrency(t *testing.T) {
	database, _ := db.InitTestDB()
	room := model.Room{
		Name:              "Single Penthouse",
		Type:              "penthouse",
		BasePricePerNight: 5000000,
		CapacityGuests:    2,
		TotalUnits:        1,
		IsActive:          true,
	}
	database.Create(&room)

	checkIn := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	checkOut := time.Now().Add(48 * time.Hour).Format("2006-01-02")

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			req := service.HoldBookingRequest{
				GuestName:       "Concurrent Guest",
				GuestPhone:      "08123456789",
				RoomID:          room.ID,
				CheckInDate:     checkIn,
				CheckOutDate:    checkOut,
				TotalRooms:      1,
				HoldDurationMin: 120,
			}
			_, err := service.HoldBooking(database, req)
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 booking to succeed for 1 unit, got %d", successCount)
	}
}
