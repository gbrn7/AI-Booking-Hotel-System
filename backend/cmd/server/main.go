package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raygbrn/ai-booking-hotel-system/internal/agent"
	"github.com/raygbrn/ai-booking-hotel-system/internal/buffer"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/handler"
	"github.com/raygbrn/ai-booking-hotel-system/internal/memory"
	"github.com/raygbrn/ai-booking-hotel-system/internal/service"
	"github.com/raygbrn/ai-booking-hotel-system/internal/tools"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting Hotel AI Booking Backend on port %s...", cfg.Port)

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	if err := db.SeedDefaultData(database); err != nil {
		log.Printf("Warning: Seeding default data encountered error: %v", err)
	}

	// Start background worker for hold expiration
	expiryInterval := time.Duration(cfg.ExpiryCheckIntervalSeconds) * time.Second
	stopWorker := service.StartHoldExpiryWorker(database, expiryInterval)
	defer close(stopWorker)

	// Core domain services
	ts := tools.NewToolService(database)
	mem := memory.NewMemory(database)
	takeover := memory.NewTakeoverManager(database)
	waClient := whatsapp.NewWhatsAppClient(cfg)

	ag, err := agent.NewAgent(cfg, ts, cfg.SystemPromptPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize AI agent (%v). Falling back to mock agent.", err)
	}

	debouncer := buffer.NewDebounceBuffer(cfg, mem, takeover, ag, waClient)

	// HTTP Handlers
	roomH := handler.NewRoomHandler(database)
	chatH := handler.NewChatHandler(database, mem, takeover, ag)
	webhookH := handler.NewWebhookHandler(debouncer)
	adminH := handler.NewAdminHandler(database, cfg, mem, takeover, waClient)

	r := gin.Default()
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// Public APIs
	api := r.Group("/api")
	{
		api.GET("/rooms", roomH.GetRooms)
		api.GET("/availability", roomH.GetRoomAvailability)
		api.POST("/bookings/hold", roomH.DirectHoldBooking)
		api.GET("/bookings/:code", roomH.GetBookingByCode)

		// Guest Web Chat SSE
		api.POST("/chat/stream", chatH.HandleChatStream)
		api.GET("/chat/history", chatH.HandleGetChatHistory)

		// Admin & Receptionist Dashboard
		admin := api.Group("/admin")
		{
			admin.POST("/login", adminH.Login)
			admin.GET("/stats", adminH.GetStats)
			admin.GET("/conversations", adminH.GetConversations)
			admin.GET("/conversations/:id/messages", adminH.GetConversationMessages)
			admin.POST("/conversations/:id/takeover", adminH.ToggleTakeover)
			admin.POST("/conversations/:id/reply", adminH.ReplyManual)
			admin.GET("/bookings", adminH.GetBookings)
			admin.POST("/bookings/:id/confirm", adminH.ConfirmBooking)
		}
	}

	// Evolution API WhatsApp webhook
	r.POST("/webhook", webhookH.HandleWhatsAppWebhook)

	log.Printf("Server running and listening at :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server shutdown with error: %v", err)
	}
}
