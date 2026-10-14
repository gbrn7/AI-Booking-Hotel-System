package memory

import (
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

type Memory struct {
	db *gorm.DB
}

func NewMemory(db *gorm.DB) *Memory {
	return &Memory{db: db}
}

// GetOrCreateConversation retrieves existing session or creates a new one.
func (m *Memory) GetOrCreateConversation(sessionID, channel string) (*model.Conversation, error) {
	var conv model.Conversation
	err := m.db.Where("id = ?", sessionID).First(&conv).Error
	if err == nil {
		return &conv, nil
	}

	conv = model.Conversation{
		ID:              sessionID,
		Channel:         channel,
		IsHumanTakeover: false,
		LastMessageAt:   time.Now(),
		CreatedAt:       time.Now(),
	}
	if err := m.db.Create(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

// SaveMessage appends a message record to conversation history.
func (m *Memory) SaveMessage(sessionID, sender, content, thought string) error {
	msg := model.Message{
		ConversationID: sessionID,
		Sender:         sender,
		Content:        content,
		ThoughtProcess: thought,
		CreatedAt:      time.Now(),
	}
	if err := m.db.Create(&msg).Error; err != nil {
		return err
	}

	return m.db.Model(&model.Conversation{}).
		Where("id = ?", sessionID).
		Update("last_message_at", time.Now()).Error
}

// GetRecentHistory retrieves the latest `limit` messages (in chronological order) as GenAI Content items.
func (m *Memory) GetRecentHistory(sessionID string, limit int) ([]*genai.Content, error) {
	var msgs []model.Message
	err := m.db.Where("conversation_id = ?", sessionID).
		Order("created_at desc").
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}

	// Restore chronological order (oldest first).
	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	// History should start with a user turn.
	for len(msgs) > 0 && msgs[0].Sender != model.SenderUser {
		msgs = msgs[1:]
	}

	var contents []*genai.Content
	for _, msg := range msgs {
		role := "user"
		if msg.Sender == model.SenderAssistant || msg.Sender == model.SenderHumanAgent {
			role = "model"
		}
		contents = append(contents, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}
	return contents, nil
}
