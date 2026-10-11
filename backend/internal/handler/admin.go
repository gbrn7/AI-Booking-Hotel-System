package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/memory"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db       *gorm.DB
	cfg      *config.Config
	mem      *memory.Memory
	takeover *memory.TakeoverManager
	waClient *whatsapp.WhatsAppClient
}

func NewAdminHandler(db *gorm.DB, cfg *config.Config, mem *memory.Memory, to *memory.TakeoverManager, wa *whatsapp.WhatsAppClient) *AdminHandler {
	return &AdminHandler{
		db:       db,
		cfg:      cfg,
		mem:      mem,
		takeover: to,
		waClient: wa,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates staff users and returns a JWT token.
func (h *AdminHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username dan password wajib diisi"})
		return
	}

	var user model.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "kredensial tidak valid"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "kredensial tidak valid"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"name":     user.Name,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.cfg.NextAuthSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
			"role":     user.Role,
		},
	})
}

// GetConversations lists all active customer chats.
func (h *AdminHandler) GetConversations(c *gin.Context) {
	var convs []model.Conversation
	if err := h.db.Preload("Guest").Order("last_message_at desc").Find(&convs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"conversations": convs})
}

// GetConversationMessages lists transcript history for a session.
func (h *AdminHandler) GetConversationMessages(c *gin.Context) {
	id := c.Param("id")
	var msgs []model.Message
	if err := h.db.Where("conversation_id = ?", id).Order("created_at asc").Find(&msgs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": msgs})
}

type TakeoverToggleRequest struct {
	IsTakeover bool `json:"is_takeover"`
}

// ToggleTakeover enables or disables AI responses for a conversation.
func (h *AdminHandler) ToggleTakeover(c *gin.Context) {
	id := c.Param("id")
	var req TakeoverToggleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.takeover.SetTakeover(id, req.IsTakeover); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":        id,
		"is_human_takeover": req.IsTakeover,
		"message":           "Status takeover berhasil diubah",
	})
}

type ManualReplyRequest struct {
	Message string `json:"message" binding:"required"`
}

// ReplyManual sends a staff reply directly into the conversation.
func (h *AdminHandler) ReplyManual(c *gin.Context) {
	id := c.Param("id")
	var req ManualReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = h.mem.SaveMessage(id, model.SenderHumanAgent, req.Message, "Manual staff response")

	// If WhatsApp channel, forward through Evolution API
	if strings.HasPrefix(id, "wa:") {
		phone := strings.TrimPrefix(id, "wa:")
		if h.waClient != nil {
			_ = h.waClient.SendText(phone, req.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent", "message": req.Message})
}

// GetBookings retrieves reservations with optional status filter.
func (h *AdminHandler) GetBookings(c *gin.Context) {
	status := c.Query("status")
	q := h.db.Preload("Room").Preload("Guest").Model(&model.Booking{})
	if status != "" {
		q = q.Where("status = ?", strings.ToUpper(status))
	}

	var bookings []model.Booking
	if err := q.Order("created_at desc").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
}

// ConfirmBooking updates a HOLD booking to CONFIRMED.
func (h *AdminHandler) ConfirmBooking(c *gin.Context) {
	id := c.Param("id")
	var booking model.Booking
	if err := h.db.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking tidak ditemukan"})
		return
	}

	now := time.Now()
	booking.Status = model.StatusConfirmed
	booking.ConfirmedAt = &now
	if err := h.db.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran diverifikasi. Booking status telah CONFIRMED.",
		"booking": booking,
	})
}

// GetStats returns summary KPI metrics for the hotel admin.
func (h *AdminHandler) GetStats(c *gin.Context) {
	var totalRooms int64
	var totalHold int64
	var totalConfirmed int64
	var activeConversations int64

	h.db.Model(&model.Room{}).Where("is_active = ?", true).Count(&totalRooms)
	h.db.Model(&model.Booking{}).Where("status = ?", model.StatusHold).Count(&totalHold)
	h.db.Model(&model.Booking{}).Where("status = ?", model.StatusConfirmed).Count(&totalConfirmed)
	h.db.Model(&model.Conversation{}).Count(&activeConversations)

	c.JSON(http.StatusOK, gin.H{
		"total_rooms":          totalRooms,
		"pending_hold_count":   totalHold,
		"confirmed_count":      totalConfirmed,
		"active_conversations": activeConversations,
	})
}
