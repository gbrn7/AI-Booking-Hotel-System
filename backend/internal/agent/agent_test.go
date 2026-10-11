package agent_test

import (
	"context"
	"testing"

	"github.com/raygbrn/ai-booking-hotel-system/internal/agent"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/tools"
)

func TestAgentInitializationAndMockFallback(t *testing.T) {
	database, _ := db.InitTestDB()
	ts := tools.NewToolService(database)

	cfg := &config.Config{
		GeminiModel:  "gemini-2.5-flash",
		GeminiAPIKey: "", // Empty to test offline mock behavior
	}

	ag, err := agent.NewAgent(cfg, ts, "../../prompts/system.txt")
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}

	reply, _, _, err := ag.RunAgent(context.Background(), nil, "Halo, apakah ada kamar kosong?", "081234567890")
	if err != nil {
		t.Fatalf("unexpected error running agent: %v", err)
	}

	if reply == "" {
		t.Errorf("expected non-empty agent reply")
	}
}
