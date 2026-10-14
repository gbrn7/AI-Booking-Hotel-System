package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raygbrn/ai-booking-hotel-system/internal/agent"
	"github.com/raygbrn/ai-booking-hotel-system/internal/memory"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
	"gorm.io/gorm"
)

type ChatHandler struct {
	db       *gorm.DB
	mem      *memory.Memory
	takeover *memory.TakeoverManager
	agent    *agent.Agent
}

func NewChatHandler(db *gorm.DB, mem *memory.Memory, to *memory.TakeoverManager, ag *agent.Agent) *ChatHandler {
	return &ChatHandler{
		db:       db,
		mem:      mem,
		takeover: to,
		agent:    ag,
	}
}

type ChatStreamRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Phone     string `json:"phone"`
}

// HandleChatStream handles real-time Server-Sent Events (SSE) for the frontend Live Chat Widget.
func (h *ChatHandler) HandleChatStream(c *gin.Context) {
	var req ChatStreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := req.SessionID
	if !strings.HasPrefix(sessionID, "web:") && !strings.HasPrefix(sessionID, "wa:") {
		sessionID = "web:" + sessionID
	}

	_, _ = h.mem.GetOrCreateConversation(sessionID, model.ChannelWeb)
	_ = h.mem.SaveMessage(sessionID, model.SenderUser, req.Message, "")

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	// If human takeover is active, do not run agent
	if h.takeover.IsTakeoverActive(sessionID) {
		msgData, _ := json.Marshal(map[string]any{
			"type":    "takeover",
			"content": "Pesan Kakak telah diterima oleh tim resepsionis kami. Staf front desk akan segera membalas di sini 😊",
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", msgData)
		flusher.Flush()
		return
	}

	history, err := h.mem.GetRecentHistory(sessionID, 20)
	if err != nil {
		history = nil
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	reply, _, thought, err := h.agent.RunAgent(ctx, history, req.Message, req.Phone)
	if err != nil {
		reply = agent.FallbackMessage
	}

	_ = h.mem.SaveMessage(sessionID, model.SenderAssistant, reply, thought)

	bubbles := whatsapp.SplitBubbles(reply)
	c.Stream(func(w io.Writer) bool {
		for i, bubble := range bubbles {
			bubbleData, _ := json.Marshal(map[string]any{
				"type":    "bubble",
				"index":   i,
				"content": bubble,
				"thought": thought,
			})
			fmt.Fprintf(w, "data: %s\n\n", bubbleData)
			flusher.Flush()
			time.Sleep(300 * time.Millisecond)
		}
		doneData, _ := json.Marshal(map[string]any{"type": "done"})
		fmt.Fprintf(w, "data: %s\n\n", doneData)
		flusher.Flush()
		return false
	})
}

// HandleGetChatHistory returns past messages for a conversation session.
func (h *ChatHandler) HandleGetChatHistory(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	if !strings.HasPrefix(sessionID, "web:") && !strings.HasPrefix(sessionID, "wa:") {
		sessionID = "web:" + sessionID
	}

	var messages []model.Message
	err := h.db.Where("conversation_id = ?", sessionID).
		Order("created_at asc").
		Find(&messages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var conv model.Conversation
	_ = h.db.Where("id = ?", sessionID).First(&conv).Error

	c.JSON(http.StatusOK, gin.H{
		"session_id":        sessionID,
		"is_human_takeover": conv.IsHumanTakeover,
		"messages":          messages,
	})
}
