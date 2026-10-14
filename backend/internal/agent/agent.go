package agent

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/tools"
	"google.golang.org/genai"
)

const (
	MaxToolIterations = 5
	FallbackMessage   = "Mohon maaf Kak, sistem concierge kami sedang mengalami lonjakan trafik🙏[NEXT]" +
		"Staf resepsionis kami akan segera membantu Kakak secara manual ya 😊"
)

var (
	HariNama  = []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	BulanNama = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli",
		"Agustus", "September", "Oktober", "November", "Desember"}
)

type Agent struct {
	client       *genai.Client
	model        string
	systemPrompt string
	tools        *tools.ToolService
	genConfig    *genai.GenerateContentConfig
}

func NewAgent(cfg *config.Config, ts *tools.ToolService, promptPath string) (*Agent, error) {
	sysPrompt, err := RenderSystemPrompt(promptPath, "Grand Oasis Hotel")
	if err != nil {
		sysPrompt = "Kamu adalah Rama, Concierge & Sales Specialist resmi Grand Oasis Hotel."
	}

	var client *genai.Client
	if cfg.GeminiAPIKey != "" {
		ctx := context.Background()
		c, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  cfg.GeminiAPIKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize genai client: %w", err)
		}
		client = c
	}

	genConfig := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: sysPrompt}},
		},
		Tools: []*genai.Tool{
			{FunctionDeclarations: ts.GetFunctionDeclarations()},
		},
		Temperature: genai.Ptr(float32(0.7)),
	}

	return &Agent{
		client:       client,
		model:        cfg.GeminiModel,
		systemPrompt: sysPrompt,
		tools:        ts,
		genConfig:    genConfig,
	}, nil
}

// RunAgent processes a user message through the Gemini agentic loop, executing tools as requested.
func (a *Agent) RunAgent(ctx context.Context, history []*genai.Content, userMessage, customerPhone string) (string, []*genai.Content, string, error) {
	if a.client == nil {
		// Mock response for development when API key is not yet set
		mockReply := "Halo Kak! Salam hangat dari Grand Oasis Hotel 🌊[NEXT]" +
			"Saya Rama, concierge digital Anda. Ada tipe kamar atau tanggal liburan yang ingin Anda tanyakan?"
		return mockReply, history, "Development mock mode (no GEMINI_API_KEY set)", nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	now := time.Now().In(loc)
	weekday := int(now.Weekday())
	dateCtx := fmt.Sprintf("Hari ini %s, %d %s %d (%s), pukul %02d:%02d WIB",
		HariNama[weekday], now.Day(), BulanNama[now.Month()-1], now.Year(),
		now.Format("2006-01-02"), now.Hour(), now.Minute())

	userContext := fmt.Sprintf("[Konteks Sistem: %s. Nomor telepon/kontak tamu: %s]\n%s", dateCtx, customerPhone, userMessage)

	// Handlers save the incoming message before loading history; drop it here to avoid sending it twice.
	if n := len(history); n > 0 {
		last := history[n-1]
		if last.Role == "user" && len(last.Parts) == 1 && last.Parts[0].Text == userMessage {
			history = history[:n-1]
		}
	}

	contents := append([]*genai.Content{}, history...)
	contents = append(contents, &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: userContext},
		},
	})

	var thoughtAccumulator strings.Builder

	for iter := 0; iter < MaxToolIterations; iter++ {
		resp, err := a.client.Models.GenerateContent(ctx, a.model, contents, a.genConfig)
		if err != nil {
			log.Printf("[Agent] Gemini API error: %v", err)
			return FallbackMessage, contents, thoughtAccumulator.String(), err
		}

		if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
			return FallbackMessage, contents, thoughtAccumulator.String(), nil
		}

		cand := resp.Candidates[0]
		contents = append(contents, cand.Content)

		functionCalls := resp.FunctionCalls()
		if len(functionCalls) == 0 {
			replyText := resp.Text()
			if replyText == "" {
				replyText = FallbackMessage
			}
			return replyText, contents, thoughtAccumulator.String(), nil
		}

		var resultParts []*genai.Part
		for _, fc := range functionCalls {
			log.Printf("  [Agent Tool Call] %s(%v)", fc.Name, fc.Args)
			thoughtAccumulator.WriteString(fmt.Sprintf("Tool call: %s(%v)\n", fc.Name, fc.Args))
			result := a.tools.ExecuteTool(fc.Name, fc.Args)
			resultParts = append(resultParts, &genai.Part{
				FunctionResponse: &genai.FunctionResponse{
					Name:     fc.Name,
					Response: map[string]any{"result": result},
				},
			})
		}
		contents = append(contents, &genai.Content{Role: "user", Parts: resultParts})
	}

	return FallbackMessage, contents, thoughtAccumulator.String(), nil
}
