package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"gorm.io/gorm"
)

var (
	ErrInvalidDateRange = errors.New("tanggal check-out harus setelah tanggal check-in")
	ErrRoomNotFound     = errors.New("kamar tidak ditemukan")
	ErrNoRoomsAvailable = errors.New("kamar penuh untuk rentang tanggal yang dipilih")
	ErrInvalidGuestInfo = errors.New("nama dan nomor telepon tamu wajib diisi")
)

// Mutex registry per room to prevent concurrent overbooking
var (
	roomMutexMap   = make(map[uint]*sync.Mutex)
	roomRegistryMu sync.Mutex
)

func getRoomMutex(roomID uint) *sync.Mutex {
	roomRegistryMu.Lock()
	defer roomRegistryMu.Unlock()
	if mu, exists := roomMutexMap[roomID]; exists {
		return mu
	}
	mu := &sync.Mutex{}
	roomMutexMap[roomID] = mu
	return mu
}

// RoomAvailability contains remaining units and pricing for a date range.
type RoomAvailability struct {
	RoomID              uint     `json:"room_id"`
	Name                string   `json:"name"`
	Type                string   `json:"type"`
	BasePricePerNight   float64  `json:"base_price_per_night"`
	CapacityGuests      int      `json:"capacity_guests"`
	TotalUnits          int      `json:"total_units"`
	AvailableUnits      int      `json:"available_units"`
	TotalNights         int      `json:"total_nights"`
	EstimatedTotalPrice float64  `json:"estimated_total_price"`
	Description         string   `json:"description"`
	Amenities           string   `json:"amenities"`
	ImageURLs           string   `json:"image_urls"`
}

// HoldBookingRequest holds parameters for creating a temporary reservation.
type HoldBookingRequest struct {
	GuestName       string `json:"guest_name"`
	GuestPhone      string `json:"guest_phone"`
	GuestEmail      string `json:"guest_email"`
	RoomID          uint   `json:"room_id"`
	CheckInDate     string `json:"check_in_date"`  // YYYY-MM-DD
	CheckOutDate    string `json:"check_out_date"` // YYYY-MM-DD
	TotalRooms      int    `json:"total_rooms"`
	SpecialRequests string `json:"special_requests"`
	HoldDurationMin int    `json:"hold_duration_min"`
}

// CheckAvailability calculates night-by-night room inventory between check-in and check-out.
func CheckAvailability(db *gorm.DB, checkInStr, checkOutStr, roomType string) ([]RoomAvailability, error) {
	checkIn, err := time.Parse("2006-01-02", checkInStr)
	if err != nil {
		return nil, fmt.Errorf("format check-in tidak valid (YYYY-MM-DD): %w", err)
	}
	checkOut, err := time.Parse("2006-01-02", checkOutStr)
	if err != nil {
		return nil, fmt.Errorf("format check-out tidak valid (YYYY-MM-DD): %w", err)
	}

	if !checkOut.After(checkIn) {
		return nil, ErrInvalidDateRange
	}

	totalNights := int(checkOut.Sub(checkIn).Hours() / 24)
	if totalNights < 1 {
		totalNights = 1
	}

	query := db.Model(&model.Room{}).Where("is_active = ?", true)
	if roomType != "" {
		query = query.Where("LOWER(type) = ?", strings.ToLower(roomType))
	}

	var rooms []model.Room
	if err := query.Find(&rooms).Error; err != nil {
		return nil, err
	}

	var results []RoomAvailability
	for _, room := range rooms {
		availableUnits, err := calculateAvailableUnitsForRoom(db, room.ID, room.TotalUnits, checkIn, totalNights)
		if err != nil {
			return nil, err
		}

		results = append(results, RoomAvailability{
			RoomID:              room.ID,
			Name:                room.Name,
			Type:                room.Type,
			BasePricePerNight:   room.BasePricePerNight,
			CapacityGuests:      room.CapacityGuests,
			TotalUnits:          room.TotalUnits,
			AvailableUnits:      availableUnits,
			TotalNights:         totalNights,
			EstimatedTotalPrice: float64(totalNights) * room.BasePricePerNight,
			Description:         room.Description,
			Amenities:           room.Amenities,
			ImageURLs:           room.ImageURLs,
		})
	}

	return results, nil
}

// calculateAvailableUnitsForRoom checks every night in the date range.
func calculateAvailableUnitsForRoom(db *gorm.DB, roomID uint, totalUnits int, checkIn time.Time, totalNights int) (int, error) {
	maxBooked := 0

	for night := 0; night < totalNights; night++ {
		nightDate := checkIn.AddDate(0, 0, night).Format("2006-01-02")

		var bookedCount int64
		err := db.Model(&model.Booking{}).
			Where("room_id = ? AND status IN (?, ?, ?) AND check_in_date <= ? AND check_out_date > ?",
				roomID, model.StatusHold, model.StatusConfirmed, model.StatusCheckedIn, nightDate, nightDate).
			Select("COALESCE(SUM(total_rooms), 0)").
			Scan(&bookedCount).Error
		if err != nil {
			return 0, err
		}

		if int(bookedCount) > maxBooked {
			maxBooked = int(bookedCount)
		}
	}

	available := totalUnits - maxBooked
	if available < 0 {
		available = 0
	}
	return available, nil
}

// HoldBooking reserves room units under HOLD status with concurrency protection.
func HoldBooking(db *gorm.DB, req HoldBookingRequest) (*model.Booking, error) {
	if req.GuestName == "" || req.GuestPhone == "" {
		return nil, ErrInvalidGuestInfo
	}
	if req.TotalRooms <= 0 {
		req.TotalRooms = 1
	}
	if req.HoldDurationMin <= 0 {
		req.HoldDurationMin = 120
	}

	checkIn, err := time.Parse("2006-01-02", req.CheckInDate)
	if err != nil {
		return nil, fmt.Errorf("format check-in tidak valid (YYYY-MM-DD): %w", err)
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOutDate)
	if err != nil {
		return nil, fmt.Errorf("format check-out tidak valid (YYYY-MM-DD): %w", err)
	}
	if !checkOut.After(checkIn) {
		return nil, ErrInvalidDateRange
	}
	totalNights := int(checkOut.Sub(checkIn).Hours() / 24)
	if totalNights < 1 {
		totalNights = 1
	}

	// Lock mutex for this room to prevent concurrent overbooking
	mu := getRoomMutex(req.RoomID)
	mu.Lock()
	defer mu.Unlock()

	var booking *model.Booking
	err = db.Transaction(func(tx *gorm.DB) error {
		var room model.Room
		if err := tx.First(&room, req.RoomID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRoomNotFound
			}
			return err
		}

		availableUnits, err := calculateAvailableUnitsForRoom(tx, room.ID, room.TotalUnits, checkIn, totalNights)
		if err != nil {
			return err
		}

		if availableUnits < req.TotalRooms {
			return ErrNoRoomsAvailable
		}

		// Find or create guest
		var guest model.Guest
		err = tx.Where("phone_number = ?", req.GuestPhone).First(&guest).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				guest = model.Guest{
					Name:        req.GuestName,
					PhoneNumber: req.GuestPhone,
					Email:       req.GuestEmail,
				}
				if err := tx.Create(&guest).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else if guest.Name != req.GuestName && req.GuestName != "" {
			guest.Name = req.GuestName
			tx.Save(&guest)
		}

		totalPrice := float64(totalNights) * room.BasePricePerNight * float64(req.TotalRooms)
		bookingCode := generateBookingCode()
		expiresAt := time.Now().Add(time.Duration(req.HoldDurationMin) * time.Minute)

		newBooking := model.Booking{
			BookingCode:     bookingCode,
			GuestID:         guest.ID,
			RoomID:          room.ID,
			CheckInDate:     req.CheckInDate,
			CheckOutDate:    req.CheckOutDate,
			TotalNights:     totalNights,
			TotalRooms:      req.TotalRooms,
			TotalPrice:      totalPrice,
			Status:          model.StatusHold,
			HoldExpiresAt:   &expiresAt,
			SpecialRequests: req.SpecialRequests,
		}

		if err := tx.Create(&newBooking).Error; err != nil {
			return err
		}

		booking = &newBooking
		booking.Room = room
		booking.Guest = guest
		return nil
	})

	if err != nil {
		return nil, err
	}
	return booking, nil
}

func generateBookingCode() string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 5)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	now := time.Now().Format("060102")
	return fmt.Sprintf("HTL-%s-%s", now, string(b))
}
