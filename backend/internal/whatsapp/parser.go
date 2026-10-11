package whatsapp

import (
	"strings"
)

type IncomingMessage struct {
	Phone     string
	Name      string
	Text      string
	MediaType string
	Mime      string
	RawData   map[string]any
}

// ExtractIncoming parses webhook payloads from Evolution API v2.
func ExtractIncoming(payload map[string]any) *IncomingMessage {
	data, _ := payload["data"].(map[string]any)
	if data == nil {
		return nil
	}
	key, _ := data["key"].(map[string]any)
	if key == nil {
		return nil
	}
	if fromMe, _ := key["fromMe"].(bool); fromMe {
		return nil
	}

	remoteJid, _ := key["remoteJid"].(string)
	if strings.Contains(remoteJid, "@g.us") || !strings.HasSuffix(remoteJid, "@s.whatsapp.net") {
		return nil
	}

	phone := strings.Split(remoteJid, "@")[0]
	name, _ := data["pushName"].(string)
	message, _ := data["message"].(map[string]any)
	if message == nil {
		return nil
	}

	if imgMsg, ok := message["imageMessage"].(map[string]any); ok {
		caption, _ := imgMsg["caption"].(string)
		mime, _ := imgMsg["mimetype"].(string)
		if mime == "" {
			mime = "image/jpeg"
		}
		return &IncomingMessage{
			Phone:     phone,
			Name:      name,
			Text:      caption,
			MediaType: "image",
			Mime:      mime,
			RawData:   data,
		}
	}

	text, _ := message["conversation"].(string)
	if text == "" {
		if extMsg, ok := message["extendedTextMessage"].(map[string]any); ok {
			text, _ = extMsg["text"].(string)
		}
	}

	if text == "" {
		return nil
	}

	return &IncomingMessage{
		Phone:   phone,
		Name:    name,
		Text:    text,
		RawData: data,
	}
}

// SplitBubbles separates text containing [NEXT] into distinct chat bubbles.
func SplitBubbles(reply string) []string {
	parts := strings.Split(reply, "[NEXT]")
	var bubbles []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			bubbles = append(bubbles, trimmed)
		}
	}
	return bubbles
}
