package whatsapp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
)

func TestSplitBubbles(t *testing.T) {
	input := "Halo selamat datang di Grand Luxury Hotel! [NEXT] Ada yang bisa kami bantu? [NEXT] "
	bubbles := SplitBubbles(input)

	if len(bubbles) != 2 {
		t.Fatalf("expected 2 bubbles, got %d", len(bubbles))
	}
	if bubbles[0] != "Halo selamat datang di Grand Luxury Hotel!" {
		t.Errorf("unexpected bubble 0: %q", bubbles[0])
	}
	if bubbles[1] != "Ada yang bisa kami bantu?" {
		t.Errorf("unexpected bubble 1: %q", bubbles[1])
	}

	// Test lowercase [next] and mixed case
	inputLower := "Bubble satu\n\n[next]\n\nBubble dua\n[Next]\nBubble tiga"
	bubblesLower := SplitBubbles(inputLower)
	if len(bubblesLower) != 3 {
		t.Fatalf("expected 3 bubbles for case-insensitive test, got %d", len(bubblesLower))
	}
	if bubblesLower[0] != "Bubble satu" || bubblesLower[1] != "Bubble dua" || bubblesLower[2] != "Bubble tiga" {
		t.Errorf("unexpected bubbles: %v", bubblesLower)
	}
}

func TestCalculateTypingDuration(t *testing.T) {
	shortDur := calculateTypingDuration("Hi")
	if shortDur < time.Second {
		t.Errorf("expected typing duration at least 1s, got %v", shortDur)
	}

	longText := string(make([]byte, 500))
	longDur := calculateTypingDuration(longText)
	if longDur > 5*time.Second {
		t.Errorf("expected typing duration max 5s, got %v", longDur)
	}
}

func TestWhatsAppClientSendText(t *testing.T) {
	var requestedPath string
	var requestedApiKey string
	var receivedBody map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		requestedApiKey = r.Header.Get("apikey")
		_ = json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"SUCCESS"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		EvolutionAPIURL:       server.URL,
		EvolutionAPIKey:       "test-secret-key",
		EvolutionInstanceName: "raygbrn",
	}

	client := NewWhatsAppClient(cfg)
	err := client.SendText("628123456789", "Halo test pesan")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedPath := "/message/sendText/raygbrn"
	if requestedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, requestedPath)
	}
	if requestedApiKey != "test-secret-key" {
		t.Errorf("expected apikey test-secret-key, got %s", requestedApiKey)
	}
	if receivedBody["number"] != "628123456789" || receivedBody["text"] != "Halo test pesan" {
		t.Errorf("unexpected body: %+v", receivedBody)
	}
}

func TestWhatsAppClientSendText404Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Instance not found"}`))
	}))
	defer server.Close()

	cfg := &config.Config{
		EvolutionAPIURL:       server.URL,
		EvolutionAPIKey:       "test-secret-key",
		EvolutionInstanceName: "invalid-instance",
	}

	client := NewWhatsAppClient(cfg)
	err := client.SendText("628123456789", "Halo")
	if err == nil {
		t.Fatal("expected error on 404, got nil")
	}
	if !contains(err.Error(), "404") || !contains(err.Error(), "Instance not found") {
		t.Errorf("expected 404 and error body in message, got: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
