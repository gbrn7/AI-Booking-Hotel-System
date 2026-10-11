package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
)

const (
	TypingSecondsPerChar = 0.03
	TypingMinSeconds     = 1.0
	TypingMaxSeconds     = 5.0
)

type WhatsAppClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewWhatsAppClient(cfg *config.Config) *WhatsAppClient {
	return &WhatsAppClient{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func calculateTypingDuration(text string) time.Duration {
	secs := float64(len(text)) * TypingSecondsPerChar
	if secs < TypingMinSeconds {
		secs = TypingMinSeconds
	}
	if secs > TypingMaxSeconds {
		secs = TypingMaxSeconds
	}
	return time.Duration(secs * float64(time.Second))
}

// SendTyping sends presence composing indicator to WhatsApp.
func (c *WhatsAppClient) SendTyping(phone string, duration time.Duration) {
	baseURL := strings.TrimRight(c.cfg.EvolutionAPIURL, "/")
	if baseURL == "" || c.cfg.EvolutionAPIKey == "" {
		return
	}

	url := fmt.Sprintf("%s/chat/sendPresence/%s", baseURL, c.cfg.EvolutionInstanceName)
	payload := map[string]any{
		"number":   phone,
		"presence": "composing",
		"delay":    int(duration.Milliseconds()),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("apikey", c.cfg.EvolutionAPIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err == nil {
		_ = resp.Body.Close()
	}
}

// SendText sends a single text message through Evolution API.
func (c *WhatsAppClient) SendText(phone, text string) error {
	baseURL := strings.TrimRight(c.cfg.EvolutionAPIURL, "/")
	if baseURL == "" || c.cfg.EvolutionAPIKey == "" {
		log.Printf("[WhatsApp Client Mock] To: %s | Text: %s", phone, text)
		return nil
	}

	url := fmt.Sprintf("%s/message/sendText/%s", baseURL, c.cfg.EvolutionInstanceName)
	payload := map[string]any{
		"number": phone,
		"text":   text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.cfg.EvolutionAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("evolution api send text error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("evolution api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

// SendReply sends a reply split across multiple [NEXT] bubbles with realistic typing delays.
func (c *WhatsAppClient) SendReply(phone, reply string) {
	bubbles := SplitBubbles(reply)
	for i, bubble := range bubbles {
		typingDur := calculateTypingDuration(bubble)
		c.SendTyping(phone, typingDur)
		time.Sleep(typingDur)

		if err := c.SendText(phone, bubble); err != nil {
			log.Printf("[WhatsApp Client] Error sending bubble to %s: %v", phone, err)
		}

		if i < len(bubbles)-1 {
			time.Sleep(800 * time.Millisecond)
		}
	}
}
