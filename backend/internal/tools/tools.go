package tools

import (
	"fmt"
	"strings"

	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

type ToolService struct {
	db *gorm.DB
}

func NewToolService(db *gorm.DB) *ToolService {
	return &ToolService{db: db}
}

// GetFunctionDeclarations returns Gemini function schemas for all 6 hotel business tools.
func (ts *ToolService) GetFunctionDeclarations() []*genai.FunctionDeclaration {
	return []*genai.FunctionDeclaration{
		{
			Name: "search_rooms",
			Description: "Cari daftar tipe kamar hotel berdasarkan kata kunci fasilitas, tipe (standard, deluxe, suite, family_suite), " +
				"atau batas anggaran harga maksimal per malam.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"query": {
						Type:        genai.TypeString,
						Description: "Kata kunci fasilitas atau pemandangan, contoh: 'ocean view', 'bathtub', 'kolam', 'balkon'",
					},
					"room_type": {
						Type:        genai.TypeString,
						Description: "Tipe kamar spesifik: 'standard', 'deluxe', 'suite', atau 'family_suite'",
					},
					"max_price": {
						Type:        genai.TypeNumber,
						Description: "Batas harga per malam maksimal dalam Rupiah, contoh: 2000000",
					},
				},
			},
		},
		{
			Name: "check_room_availability",
			Description: "Cek ketersediaan sisa unit kamar pada rentang tanggal menginap (check-in s.d. check-out). " +
				"Wajib dipanggil saat tamu menyebutkan tanggal atau menanyakan ketersediaan kamar.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"check_in_date": {
						Type:        genai.TypeString,
						Description: "Tanggal check-in (format YYYY-MM-DD), contoh: '2026-11-01'",
					},
					"check_out_date": {
						Type:        genai.TypeString,
						Description: "Tanggal check-out (format YYYY-MM-DD), contoh: '2026-11-03'",
					},
					"room_type": {
						Type:        genai.TypeString,
						Description: "Filter tipe kamar opsional, contoh: 'deluxe'",
					},
				},
				Required: []string{"check_in_date", "check_out_date"},
			},
		},
		{
			Name: "hold_booking",
			Description: "Membuat reservasi sementara kamar berstatus HOLD (berlaku 2 jam) dengan proteksi inventori atomik. " +
				"Panggil HANYA setelah tanggal, kamar, nama, dan nomor telepon tamu sudah jelas dikonfirmasi.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"guest_name": {
						Type:        genai.TypeString,
						Description: "Nama lengkap tamu",
					},
					"guest_phone": {
						Type:        genai.TypeString,
						Description: "Nomor kontak / WhatsApp tamu",
					},
					"guest_email": {
						Type:        genai.TypeString,
						Description: "Email tamu (opsional)",
					},
					"room_id": {
						Type:        genai.TypeInteger,
						Description: "ID unik tipe kamar (didapat dari search_rooms atau check_room_availability)",
					},
					"check_in_date": {
						Type:        genai.TypeString,
						Description: "Tanggal check-in format YYYY-MM-DD",
					},
					"check_out_date": {
						Type:        genai.TypeString,
						Description: "Tanggal check-out format YYYY-MM-DD",
					},
					"total_rooms": {
						Type:        genai.TypeInteger,
						Description: "Jumlah kamar yang dipesan (default 1)",
					},
					"special_requests": {
						Type:        genai.TypeString,
						Description: "Permintaan khusus tamu, contoh: 'Lantai tinggi', 'Ranjang non-smoking'",
					},
				},
				Required: []string{"guest_name", "guest_phone", "room_id", "check_in_date", "check_out_date"},
			},
		},
		{
			Name: "get_guest_bookings",
			Description: "Mengambil data reservasi aktif atau riwayat pemesanan tamu berdasarkan nomor telepon atau kode booking.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"phone": {
						Type:        genai.TypeString,
						Description: "Nomor telepon tamu",
					},
					"booking_code": {
						Type:        genai.TypeString,
						Description: "Kode booking unik hotel, contoh: 'HTL-261011-8X9Y'",
					},
				},
			},
		},
		{
			Name: "cancel_booking",
			Description: "Membatalkan reservasi berstatus HOLD atau melepaskan unit kamar kembali ke inventori.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"booking_code": {
						Type:        genai.TypeString,
						Description: "Kode booking unik hotel",
					},
					"phone": {
						Type:        genai.TypeString,
						Description: "Nomor telepon pemesan untuk verifikasi",
					},
					"reason": {
						Type:        genai.TypeString,
						Description: "Alasan pembatalan",
					},
				},
				Required: []string{"booking_code", "phone"},
			},
		},
		{
			Name: "transfer_to_receptionist",
			Description: "Mengalihkan percakapan ke resepsionis manusia (takeover mode) saat tamu memerlukan bantuan khusus atau komplain.",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"reason": {
						Type:        genai.TypeString,
						Description: "Alasan pengalihan ke staf manusia",
					},
					"phone": {
						Type:        genai.TypeString,
						Description: "Nomor telepon tamu",
					},
				},
				Required: []string{"reason"},
			},
		},
	}
}

// ExecuteTool dispatches a function call by name with parsed argument map.
func (ts *ToolService) ExecuteTool(name string, args map[string]any) any {
	switch name {
	case "search_rooms":
		query, _ := args["query"].(string)
		roomType, _ := args["room_type"].(string)
		maxPrice, _ := args["max_price"].(float64)
		rooms, err := ts.SearchRooms(query, roomType, maxPrice)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"rooms": rooms, "count": len(rooms)}

	case "check_room_availability":
		checkIn, _ := args["check_in_date"].(string)
		checkOut, _ := args["check_out_date"].(string)
		roomType, _ := args["room_type"].(string)
		avails, err := service.CheckAvailability(ts.db, checkIn, checkOut, roomType)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"availability": avails}

	case "hold_booking":
		guestName, _ := args["guest_name"].(string)
		guestPhone, _ := args["guest_phone"].(string)
		guestEmail, _ := args["guest_email"].(string)
		checkIn, _ := args["check_in_date"].(string)
		checkOut, _ := args["check_out_date"].(string)
		specialRequests, _ := args["special_requests"].(string)

		var roomID uint
		if rVal, ok := args["room_id"].(float64); ok {
			roomID = uint(rVal)
		} else if rVal, ok := args["room_id"].(int); ok {
			roomID = uint(rVal)
		}

		totalRooms := 1
		if trVal, ok := args["total_rooms"].(float64); ok && trVal > 0 {
			totalRooms = int(trVal)
		}

		booking, err := service.HoldBooking(ts.db, service.HoldBookingRequest{
			GuestName:       guestName,
			GuestPhone:      guestPhone,
			GuestEmail:      guestEmail,
			RoomID:          roomID,
			CheckInDate:     checkIn,
			CheckOutDate:    checkOut,
			TotalRooms:      totalRooms,
			SpecialRequests: specialRequests,
			HoldDurationMin: 120,
		})
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{
			"status":            "HOLD_SUCCESS",
			"booking_code":      booking.BookingCode,
			"room_name":         booking.Room.Name,
			"check_in":          booking.CheckInDate,
			"check_out":         booking.CheckOutDate,
			"total_nights":      booking.TotalNights,
			"total_price":       booking.TotalPrice,
			"hold_expires_at":   booking.HoldExpiresAt,
			"payment_recipient": "BCA Virtual Account: 8870-1234-5678-9000 a/n Grand Oasis Hotel",
		}

	case "get_guest_bookings":
		phone, _ := args["phone"].(string)
		code, _ := args["booking_code"].(string)
		bookings, err := ts.GetGuestBookings(phone, code)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"bookings": bookings}

	case "cancel_booking":
		code, _ := args["booking_code"].(string)
		phone, _ := args["phone"].(string)
		reason, _ := args["reason"].(string)
		msg, err := ts.CancelBooking(code, phone, reason)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"result": msg}

	case "transfer_to_receptionist":
		reason, _ := args["reason"].(string)
		phone, _ := args["phone"].(string)
		msg, err := ts.TransferToReceptionist(phone, reason)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}
		return map[string]any{"result": msg}

	default:
		return map[string]any{"error": fmt.Sprintf("tool '%s' tidak dikenali", name)}
	}
}

func (ts *ToolService) SearchRooms(query, roomType string, maxPrice float64) ([]model.Room, error) {
	q := ts.db.Model(&model.Room{}).Where("is_active = ?", true)

	if roomType != "" {
		q = q.Where("LOWER(type) = ?", strings.ToLower(roomType))
	}
	if maxPrice > 0 {
		q = q.Where("base_price_per_night <= ?", maxPrice)
	}
	if query != "" {
		kw := "%" + strings.ToLower(query) + "%"
		q = q.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(amenities) LIKE ?", kw, kw, kw)
	}

	var rooms []model.Room
	if err := q.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (ts *ToolService) GetGuestBookings(phone, bookingCode string) ([]model.Booking, error) {
	q := ts.db.Preload("Room").Preload("Guest").Model(&model.Booking{})
	if bookingCode != "" {
		q = q.Where("booking_code = ?", bookingCode)
	} else if phone != "" {
		q = q.Joins("JOIN guests ON guests.id = bookings.guest_id").Where("guests.phone_number = ?", phone)
	} else {
		return nil, fmt.Errorf("masukkan nomor telepon atau kode booking")
	}

	var bookings []model.Booking
	if err := q.Order("bookings.created_at desc").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (ts *ToolService) CancelBooking(bookingCode, phone, reason string) (string, error) {
	var booking model.Booking
	err := ts.db.Preload("Guest").Where("booking_code = ?", bookingCode).First(&booking).Error
	if err != nil {
		return "", fmt.Errorf("reservasi dengan kode %s tidak ditemukan", bookingCode)
	}

	if booking.Guest.PhoneNumber != phone && phone != "" {
		return "", fmt.Errorf("nomor telepon tidak sesuai dengan data pemesan")
	}

	if booking.Status != model.StatusHold {
		return "", fmt.Errorf("reservasi berstatus %s tidak dapat dibatalkan otomatis melalui AI. Silakan hubungi resepsionis.", booking.Status)
	}

	booking.Status = model.StatusCancelled
	if err := ts.db.Save(&booking).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("Reservasi %s berhasil dibatalkan. Kuota kamar telah dilepas kembali.", bookingCode), nil
}

func (ts *ToolService) TransferToReceptionist(phone, reason string) (string, error) {
	if phone != "" {
		conversationID := "wa:" + phone
		ts.db.Model(&model.Conversation{}).
			Where("id = ? OR id = ?", conversationID, phone).
			Update("is_human_takeover", true)
	}
	return fmt.Sprintf("Percakapan telah dialihkan ke resepsionis manusia. Alasan: %s", reason), nil
}
