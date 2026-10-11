package tools_test

import (
	"testing"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/tools"
)

func TestHotelToolsExecution(t *testing.T) {
	database, err := db.InitTestDB()
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}
	if err := db.SeedDefaultData(database); err != nil {
		t.Fatalf("failed to seed data: %v", err)
	}

	ts := tools.NewToolService(database)

	// 1. Test search_rooms
	resSearch := ts.ExecuteTool("search_rooms", map[string]any{
		"query": "ocean",
	})
	searchMap, ok := resSearch.(map[string]any)
	if !ok || searchMap["count"].(int) < 1 {
		t.Fatalf("expected at least 1 room with 'ocean', got %+v", resSearch)
	}

	// 2. Test check_room_availability
	checkIn := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	checkOut := time.Now().Add(72 * time.Hour).Format("2006-01-02")

	resAvail := ts.ExecuteTool("check_room_availability", map[string]any{
		"check_in_date":  checkIn,
		"check_out_date": checkOut,
		"room_type":      "deluxe",
	})
	availMap, ok := resAvail.(map[string]any)
	if !ok || availMap["error"] != nil {
		t.Fatalf("unexpected error in check_room_availability: %+v", resAvail)
	}

	// 3. Test hold_booking
	resHold := ts.ExecuteTool("hold_booking", map[string]any{
		"guest_name":       "Budi Santoso",
		"guest_phone":      "081234567890",
		"room_id":          float64(2), // Deluxe room
		"check_in_date":    checkIn,
		"check_out_date":   checkOut,
		"total_rooms":      float64(1),
		"special_requests": "Pillow extra",
	})
	holdMap, ok := resHold.(map[string]any)
	if !ok || holdMap["status"] != "HOLD_SUCCESS" {
		t.Fatalf("expected HOLD_SUCCESS, got %+v", resHold)
	}
	bookingCode, _ := holdMap["booking_code"].(string)

	// 4. Test get_guest_bookings
	resGet := ts.ExecuteTool("get_guest_bookings", map[string]any{
		"phone": "081234567890",
	})
	getMap, ok := resGet.(map[string]any)
	if !ok || getMap["error"] != nil {
		t.Fatalf("expected bookings list, got %+v", resGet)
	}

	// 5. Test cancel_booking
	resCancel := ts.ExecuteTool("cancel_booking", map[string]any{
		"booking_code": bookingCode,
		"phone":        "081234567890",
		"reason":       "Perubahan jadwal rencana",
	})
	cancelMap, ok := resCancel.(map[string]any)
	if !ok || cancelMap["error"] != nil {
		t.Fatalf("expected cancel success, got %+v", resCancel)
	}
}
