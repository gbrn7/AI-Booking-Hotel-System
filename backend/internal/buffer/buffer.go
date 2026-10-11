package buffer

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/agent"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/memory"
	"github.com/raygbrn/ai-booking-hotel-system/internal/model"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
	"google.golang.org/genai"
)

type AgentRunner interface {
	RunAgent(ctx context.Context, history []*genai.Content, userMessage, customerPhone string) (string, []*genai.Content, string, error)
}

type WhatsAppSender interface {
	SendReply(phone, reply string)
}

type DebounceBuffer struct {
	cfg         *config.Config
	mem         *memory.Memory
	takeover    *memory.TakeoverManager
	agentRunner AgentRunner
	waSender    WhatsAppSender
	buffers     map[string][]*whatsapp.IncomingMessage
	timers      map[string]*time.Timer
	locks       map[string]*sync.Mutex
	mapMutex    sync.Mutex
}

func NewDebounceBuffer(
	cfg *config.Config,
	mem *memory.Memory,
	to *memory.TakeoverManager,
	ag *agent.Agent,
	wa *whatsapp.WhatsAppClient,
) *DebounceBuffer {
	return &DebounceBuffer{
		cfg:         cfg,
		mem:         mem,
		takeover:    to,
		agentRunner: ag,
		waSender:    wa,
		buffers:     make(map[string][]*whatsapp.IncomingMessage),
		timers:      make(map[string]*time.Timer),
		locks:       make(map[string]*sync.Mutex),
	}
}

func (db *DebounceBuffer) SetAgentRunner(ar AgentRunner) {
	db.agentRunner = ar
}

func (db *DebounceBuffer) SetWhatsAppSender(ws WhatsAppSender) {
	db.waSender = ws
}

func (db *DebounceBuffer) getLock(phone string) *sync.Mutex {
	db.mapMutex.Lock()
	defer db.mapMutex.Unlock()
	if _, ok := db.locks[phone]; !ok {
		db.locks[phone] = &sync.Mutex{}
	}
	return db.locks[phone]
}

func (db *DebounceBuffer) Add(msg *whatsapp.IncomingMessage) int {
	db.mapMutex.Lock()
	db.buffers[msg.Phone] = append(db.buffers[msg.Phone], msg)
	count := len(db.buffers[msg.Phone])

	if t, exists := db.timers[msg.Phone]; exists {
		t.Stop()
	}

	phone := msg.Phone
	duration := time.Duration(db.cfg.WADebounceSeconds) * time.Second
	if duration <= 0 {
		duration = 4 * time.Second
	}

	db.timers[phone] = time.AfterFunc(duration, func() {
		db.processBufferedMessages(phone)
	})
	db.mapMutex.Unlock()

	return count
}

func (db *DebounceBuffer) processBufferedMessages(phone string) {
	mu := db.getLock(phone)
	mu.Lock()
	defer mu.Unlock()

	db.mapMutex.Lock()
	msgs := db.buffers[phone]
	delete(db.buffers, phone)
	delete(db.timers, phone)
	db.mapMutex.Unlock()

	if len(msgs) == 0 {
		return
	}

	var textParts []string
	for _, m := range msgs {
		if strings.TrimSpace(m.Text) != "" {
			textParts = append(textParts, strings.TrimSpace(m.Text))
		}
	}
	combinedText := strings.Join(textParts, "\n")
	if combinedText == "" {
		return
	}

	sessionID := "wa:" + phone
	_, _ = db.mem.GetOrCreateConversation(sessionID, model.ChannelWhatsApp)

	// Save incoming user message
	_ = db.mem.SaveMessage(sessionID, model.SenderUser, combinedText, "")

	// Check if human takeover is active
	if db.takeover.IsTakeoverActive(sessionID) {
		log.Printf("[Buffer] Takeover is ACTIVE for %s. Message logged without calling AI agent.", phone)
		return
	}

	if db.agentRunner == nil {
		log.Printf("[Buffer] No agent runner registered for %s.", phone)
		return
	}

	history, err := db.mem.GetRecentHistory(sessionID, 10)
	if err != nil {
		log.Printf("[Buffer] Failed to get conversation history: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	reply, _, thought, err := db.agentRunner.RunAgent(ctx, history, combinedText, phone)
	if err != nil {
		log.Printf("[Buffer] Agent error for %s: %v", phone, err)
	}

	_ = db.mem.SaveMessage(sessionID, model.SenderAssistant, reply, thought)

	if db.waSender != nil {
		db.waSender.SendReply(phone, reply)
	}
}
