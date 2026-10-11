package buffer_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/buffer"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/memory"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
	"google.golang.org/genai"
)

type mockAgentRunner struct {
	mu           sync.Mutex
	receivedMsgs []string
}

func (m *mockAgentRunner) RunAgent(ctx context.Context, history []*genai.Content, userMessage, customerPhone string) (string, []*genai.Content, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.receivedMsgs = append(m.receivedMsgs, userMessage)
	return "Terima kasih telah menghubungi kami![NEXT]Ada yang bisa dibantu?", nil, "", nil
}

type mockWASender struct {
	mu      sync.Mutex
	replies []string
}

func (m *mockWASender) SendReply(phone, reply string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.replies = append(m.replies, reply)
}

func TestDebounceBufferCombinesMessages(t *testing.T) {
	database, _ := db.InitTestDB()
	mem := memory.NewMemory(database)
	to := memory.NewTakeoverManager(database)

	cfg := &config.Config{
		WADebounceSeconds: 1, // 1 second for fast unit test
	}

	buf := buffer.NewDebounceBuffer(cfg, mem, to, nil, nil)
	agentMock := &mockAgentRunner{}
	waMock := &mockWASender{}
	buf.SetAgentRunner(agentMock)
	buf.SetWhatsAppSender(waMock)

	// Send 3 rapid messages
	buf.Add(&whatsapp.IncomingMessage{Phone: "0812345", Text: "Halo kak"})
	buf.Add(&whatsapp.IncomingMessage{Phone: "0812345", Text: "Mau tanya"})
	buf.Add(&whatsapp.IncomingMessage{Phone: "0812345", Text: "Kamar deluxe kosong tgl 15?"})

	// Wait for debounce timer to fire
	time.Sleep(1500 * time.Millisecond)

	agentMock.mu.Lock()
	msgCount := len(agentMock.receivedMsgs)
	var combined string
	if msgCount > 0 {
		combined = agentMock.receivedMsgs[0]
	}
	agentMock.mu.Unlock()

	if msgCount != 1 {
		t.Fatalf("expected exactly 1 combined invocation to agent runner, got %d", msgCount)
	}

	expectedContent := "Halo kak\nMau tanya\nKamar deluxe kosong tgl 15?"
	if combined != expectedContent {
		t.Errorf("expected combined text %q, got %q", expectedContent, combined)
	}

	waMock.mu.Lock()
	replyCount := len(waMock.replies)
	waMock.mu.Unlock()

	if replyCount != 1 {
		t.Errorf("expected 1 reply dispatched to whatsapp, got %d", replyCount)
	}
}
