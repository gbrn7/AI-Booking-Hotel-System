package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
	"gorm.io/gorm"
)

type RoomHandler struct {
	db *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

// GetRooms returns all active hotel rooms.
func (h *RoomHandler) GetRooms(c *gin.Context) {
	var rooms []model.Room
	if err := h.db.Where("is_active = ?", true).Order("base_price_per_night asc").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// GetRoomAvailability checks remaining units for requested dates.
func (h *RoomHandler) GetRoomAvailability(c *gin.Context) {
	checkIn := c.Query("check_in")
	checkOut := c.Query("check_out")
	roomType := c.Query("room_type")

	if checkIn == "" || checkOut == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "check_in dan check_out wajib diisi (YYYY-MM-DD)"})
		return
	}

	avails, err := service.CheckAvailability(h.db, checkIn, checkOut, roomType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availability": avails})
}

// DirectHoldBooking allows creating a hold booking from web forms.
func (h *RoomHandler) DirectHoldBooking(c *gin.Context) {
	var req service.HoldBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := service.HoldBooking(h.db, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reservasi berhasil dibuat (Status: HOLD)",
		"booking": booking,
	})
}

// GetBookingByCode retrieves booking details by booking code.
func (h *RoomHandler) GetBookingByCode(c *gin.Context) {
	code := c.Param("code")
	var booking model.Booking
	if err := h.db.Preload("Room").Preload("Guest").Where("booking_code = ?", code).First(&booking).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"booking": booking})
}
