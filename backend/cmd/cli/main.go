package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raygbrn/ai-booking-hotel-system/internal/agent"
	"github.com/raygbrn/ai-booking-hotel-system/internal/config"
	"github.com/raygbrn/ai-booking-hotel-system/internal/db"
	"github.com/raygbrn/ai-booking-hotel-system/internal/tools"
	"github.com/raygbrn/ai-booking-hotel-system/internal/whatsapp"
	"google.golang.org/genai"
)

func main() {
	cfg := config.LoadConfig()
	// Use sqlite for instant CLI testing if postgres isn't running
	cfg.DBDriver = "sqlite"

	fmt.Println("==================================================================")
	fmt.Println("  🏨 Grand Oasis Hotel — Rama AI Concierge Simulator (CLI)")
	fmt.Println("==================================================================")
	fmt.Printf("Model    : %s\n", cfg.GeminiModel)
	if cfg.GeminiAPIKey == "" {
		fmt.Println("Mode     : ⚠️  Mock Development (GEMINI_API_KEY belum diset)")
	} else {
		fmt.Println("Mode     : 🟢 Live Gemini GenAI SDK")
	}
	fmt.Println("Ketik pesan Anda dan tekan ENTER. Ketik 'exit' atau 'quit' untuk keluar.")
	fmt.Println("==================================================================")

	database, err := db.InitDB(cfg)
	if err != nil {
		fmt.Printf("Gagal koneksi database: %v\n", err)
		return
	}
	_ = db.SeedDefaultData(database)

	ts := tools.NewToolService(database)
	ag, err := agent.NewAgent(cfg, ts, cfg.SystemPromptPath)
	if err != nil {
		fmt.Printf("Gagal inisialisasi AI agent: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	var history []*genai.Content
	customerPhone := "081234567890"

	fmt.Printf("\n[Rama]: Selamat datang di Grand Oasis Hotel! Ada yang bisa Rama bantu hari ini, Kak? 😊\n\n")

	for {
		fmt.Print("Anda: ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if text == "exit" || text == "quit" {
			fmt.Println("\nTerima kasih telah berkunjung ke Grand Oasis Hotel. Sampai jumpa!")
			break
		}

		reply, newHistory, thought, err := ag.RunAgent(context.Background(), history, text, customerPhone)
		if err != nil {
			fmt.Printf("\n[Error]: %v\n\n", err)
			continue
		}
		history = newHistory

		if thought != "" {
			fmt.Printf("\n🔍 [Internal Tool Execution]:\n%s\n", strings.TrimSpace(thought))
		}

		bubbles := whatsapp.SplitBubbles(reply)
		fmt.Println()
		for _, bubble := range bubbles {
			fmt.Printf("[Rama]: %s\n\n", bubble)
		}
	}
}
