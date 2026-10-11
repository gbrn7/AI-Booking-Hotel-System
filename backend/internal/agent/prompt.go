package agent

import (
	"fmt"
	"os"
	"strings"
)

// RenderSystemPrompt reads the prompt template file and injects dynamic hotel names or parameters.
func RenderSystemPrompt(path, hotelName string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read system prompt file (%s): %w", path, err)
	}

	content := string(data)
	if hotelName != "" {
		content = strings.ReplaceAll(content, "Grand Oasis Hotel", hotelName)
	}
	return content, nil
}
