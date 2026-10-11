package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
)

type WhatsAppClient struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewWhatsAppClient(cfg *config.Config) *WhatsAppClient {
	return &WhatsAppClient{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// SendText sends a single text message through Evolution API.
func (c *WhatsAppClient) SendText(phone, text string) error {
	if c.cfg.EvolutionAPIURL == "" || c.cfg.EvolutionAPIKey == "" {
		log.Printf("[WhatsApp Client Mock] To: %s | Text: %s", phone, text)
		return nil
	}

	url := fmt.Sprintf("%s/message/sendText/%s", c.cfg.EvolutionAPIURL, c.cfg.EvolutionInstanceName)
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
		return fmt.Errorf("evolution api returned status %d", resp.StatusCode)
	}
	return nil
}

// SendReply sends a reply split across multiple [NEXT] bubbles with realistic typing delays.
func (c *WhatsAppClient) SendReply(phone, reply string) {
	bubbles := SplitBubbles(reply)
	for i, bubble := range bubbles {
		if i > 0 {
			time.Sleep(1200 * time.Millisecond)
		}
		if err := c.SendText(phone, bubble); err != nil {
			log.Printf("[WhatsApp Client] Error sending bubble to %s: %v", phone, err)
		}
	}
}
