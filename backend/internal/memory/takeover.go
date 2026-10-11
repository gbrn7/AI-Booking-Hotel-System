package memory

import (
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"gorm.io/gorm"
)

type TakeoverManager struct {
	db *gorm.DB
}

func NewTakeoverManager(db *gorm.DB) *TakeoverManager {
	return &TakeoverManager{db: db}
}

// IsTakeoverActive checks whether human takeover is enabled for a given conversation.
func (tm *TakeoverManager) IsTakeoverActive(sessionID string) bool {
	var conv model.Conversation
	if err := tm.db.Where("id = ?", sessionID).First(&conv).Error; err != nil {
		return false
	}
	return conv.IsHumanTakeover
}

// SetTakeover toggles takeover state on or off.
func (tm *TakeoverManager) SetTakeover(sessionID string, active bool) error {
	return tm.db.Model(&model.Conversation{}).
		Where("id = ?", sessionID).
		Update("is_human_takeover", active).Error
}
