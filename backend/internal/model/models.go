package model

import (
	"time"
)

// Room represents a hotel room type with capacity and total available units.
type Room struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:100;not null" json:"name"`
	Type              string    `gorm:"size:50;not null;index" json:"type"` // standard, deluxe, suite, family_suite
	BasePricePerNight float64   `gorm:"not null" json:"base_price_per_night"`
	CapacityGuests    int       `gorm:"not null;default:2" json:"capacity_guests"`
	TotalUnits        int       `gorm:"not null;default:5" json:"total_units"`
	Description       string    `gorm:"type:text" json:"description"`
	Amenities         string    `gorm:"type:text" json:"amenities"` // JSON string e.g. ["WiFi", "AC", "Bathtub"]
	ImageURLs         string    `gorm:"type:text" json:"image_urls"` // JSON string of urls
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Guest represents a hotel customer.
type Guest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	PhoneNumber string    `gorm:"size:30;not null;uniqueIndex" json:"phone_number"`
	Email       string    `gorm:"size:100" json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookingStatus enum constants
const (
	StatusHold       = "HOLD"
	StatusConfirmed  = "CONFIRMED"
	StatusCancelled  = "CANCELLED"
	StatusCheckedIn  = "CHECKED_IN"
	StatusCheckedOut = "CHECKED_OUT"
)

// Booking represents a room reservation with dates and hold status.
type Booking struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	BookingCode       string     `gorm:"size:50;not null;uniqueIndex" json:"booking_code"`
	GuestID           uint       `gorm:"not null;index" json:"guest_id"`
	Guest             Guest      `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	RoomID            uint       `gorm:"not null;index" json:"room_id"`
	Room              Room       `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	CheckInDate       string     `gorm:"size:10;not null;index" json:"check_in_date"`  // YYYY-MM-DD
	CheckOutDate      string     `gorm:"size:10;not null;index" json:"check_out_date"` // YYYY-MM-DD
	TotalNights       int        `gorm:"not null;default:1" json:"total_nights"`
	TotalRooms        int        `gorm:"not null;default:1" json:"total_rooms"`
	TotalPrice        float64    `gorm:"not null" json:"total_price"`
	Status            string     `gorm:"size:20;not null;default:'HOLD';index" json:"status"`
	HoldExpiresAt     *time.Time `gorm:"index" json:"hold_expires_at,omitempty"`
	PaymentProofURL   string     `gorm:"size:255" json:"payment_proof_url,omitempty"`
	SpecialRequests   string     `gorm:"type:text" json:"special_requests,omitempty"`
	ConfirmedByUserID *uint      `json:"confirmed_by_user_id,omitempty"`
	ConfirmedAt       *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// Channel types
const (
	ChannelWhatsApp = "whatsapp"
	ChannelWeb      = "web"
)

// Conversation tracks customer sessions and takeover states across channels.
type Conversation struct {
	ID              string    `gorm:"primaryKey;size:100" json:"id"` // wa:628... or web:uuid
	GuestID         *uint     `gorm:"index" json:"guest_id,omitempty"`
	Guest           *Guest    `gorm:"foreignKey:GuestID" json:"guest,omitempty"`
	Channel         string    `gorm:"size:20;not null;default:'web'" json:"channel"`
	IsHumanTakeover bool      `gorm:"default:false;index" json:"is_human_takeover"`
	LastMessageAt   time.Time `gorm:"index" json:"last_message_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// Sender types
const (
	SenderUser       = "user"
	SenderAssistant  = "assistant"
	SenderAgentTool  = "agent_tool"
	SenderHumanAgent = "human_agent"
)

// Message records chat history in a conversation.
type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ConversationID string    `gorm:"size:100;not null;index" json:"conversation_id"`
	Sender         string    `gorm:"size:20;not null" json:"sender"`
	Content        string    `gorm:"type:text;not null" json:"content"`
	ThoughtProcess string    `gorm:"type:text" json:"thought_process,omitempty"`
	CreatedAt      time.Time `gorm:"index" json:"created_at"`
}

// User roles
const (
	RoleAdmin        = "admin"
	RoleReceptionist = "receptionist"
)

// User represents admin or hotel staff for dashboard authentication.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"size:100;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Role         string    `gorm:"size:20;not null;default:'receptionist'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
