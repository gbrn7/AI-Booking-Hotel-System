package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raygbrn/ai-booking-hotel-system/internal/buffer"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
)

type WebhookHandler struct {
	buffer *buffer.DebounceBuffer
}

func NewWebhookHandler(buf *buffer.DebounceBuffer) *WebhookHandler {
	return &WebhookHandler{buffer: buf}
}

// HandleWhatsAppWebhook processes incoming webhook payloads from Evolution API.
func (h *WebhookHandler) HandleWhatsAppWebhook(c *gin.Context) {
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload json"})
		return
	}

	incoming := whatsapp.ExtractIncoming(payload)
	if incoming == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored", "reason": "not a customer text/image message"})
		return
	}

	log.Printf("[Webhook] Received WhatsApp message from %s (%s): %s", incoming.Phone, incoming.Name, incoming.Text)

	bufferedCount := h.buffer.Add(incoming)
	c.JSON(http.StatusOK, gin.H{
		"status":         "buffered",
		"phone":          incoming.Phone,
		"buffered_count": bufferedCount,
	})
}
