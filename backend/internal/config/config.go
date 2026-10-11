package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                     string
	Port                       string
	DBDriver                   string
	DBHost                     string
	DBPort                     string
	DBUser                     string
	DBPassword                 string
	DBName                     string
	DBSSLMode                  string
	GeminiAPIKey               string
	GeminiModel                string
	SystemPromptPath           string
	HoldDurationMinutes        int
	ExpiryCheckIntervalSeconds int
	EvolutionAPIURL            string
	EvolutionAPIKey            string
	EvolutionInstanceName      string
	WADebounceSeconds          int
	NextAuthSecret             string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	holdDuration, _ := strconv.Atoi(getEnv("HOLD_DURATION_MINUTES", "120"))
	expiryInterval, _ := strconv.Atoi(getEnv("EXPIRY_CHECK_INTERVAL_SECONDS", "60"))
	waDebounce, _ := strconv.Atoi(getEnv("WA_DEBOUNCE_SECONDS", "4"))

	return &Config{
		AppEnv:                     getEnv("APP_ENV", "development"),
		Port:                       getEnv("PORT", "8080"),
		DBDriver:                   getEnv("DB_DRIVER", "postgres"),
		DBHost:                     getEnv("DB_HOST", "localhost"),
		DBPort:                     getEnv("DB_PORT", "5432"),
		DBUser:                     getEnv("DB_USER", "hotel_admin"),
		DBPassword:                 getEnv("DB_PASSWORD", "hotel_secret_pass"),
		DBName:                     getEnv("DB_NAME", "hotel_booking_db"),
		DBSSLMode:                  getEnv("DB_SSLMODE", "disable"),
		GeminiAPIKey:               getEnv("GEMINI_API_KEY", ""),
		GeminiModel:                getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		SystemPromptPath:           getEnv("SYSTEM_PROMPT_PATH", "./prompts/system.txt"),
		HoldDurationMinutes:        holdDuration,
		ExpiryCheckIntervalSeconds: expiryInterval,
		EvolutionAPIURL:            strings.TrimRight(getEnv("EVOLUTION_API_URL", "http://localhost:8083"), "/"),
		EvolutionAPIKey:            getEnv("EVOLUTION_API_KEY", ""),
		EvolutionInstanceName:      getEnvWithFallback("EVOLUTION_INSTANCE", "EVOLUTION_INSTANCE_NAME", "hotel-sales-agent"),
		WADebounceSeconds:          waDebounce,
		NextAuthSecret:             getEnv("NEXTAUTH_SECRET", "default_secret_key_change_me"),
	}
}

func getEnvWithFallback(primaryKey, fallbackKey, defaultVal string) string {
	if val := os.Getenv(primaryKey); val != "" {
		return val
	}
	if val := os.Getenv(fallbackKey); val != "" {
		return val
	}
	return defaultVal
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
