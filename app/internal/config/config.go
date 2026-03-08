package config

import (
	"log/slog"
	"os"
)

// Config holds all startup configuration sourced from environment variables
// and CLI flags. Validated at startup so the rest of the application can
// rely on sane defaults.
type Config struct {
	Port     string
	Dev      bool
	DataDir  string
	AuthMode string // "cac" or "" (demo)

	// AI provider (at least one group must be set for AI to work)
	OpenAIKey             string
	AzureOpenAIEndpoint   string
	AzureOpenAIKey        string
	AzureOpenAIDeployment string

	// Microsoft Graph (all 3 required for Graph features)
	GraphTenantID     string
	GraphClientID     string
	GraphClientSecret string
	GraphCloud        string // "commercial", "gcc-high", "dod"
	GraphMasterCalID  string

	// Weather location (defaults to Quantico)
	WeatherLat string
	WeatherLon string
}

// Load reads configuration from environment variables, applies defaults,
// and logs warnings for features that will run in degraded mode.
func Load() Config {
	cfg := Config{
		Port:     envOr("PORT", "8080"),
		DataDir:  envOr("DATA_DIR", "data"),
		AuthMode: os.Getenv("AUTH_MODE"),

		// AI
		OpenAIKey:             os.Getenv("OPENAI_API_KEY"),
		AzureOpenAIEndpoint:   os.Getenv("AZURE_OPENAI_ENDPOINT"),
		AzureOpenAIKey:        os.Getenv("AZURE_OPENAI_KEY"),
		AzureOpenAIDeployment: os.Getenv("AZURE_OPENAI_DEPLOYMENT"),

		// Microsoft Graph
		GraphTenantID:     os.Getenv("GRAPH_TENANT_ID"),
		GraphClientID:     os.Getenv("GRAPH_CLIENT_ID"),
		GraphClientSecret: os.Getenv("GRAPH_CLIENT_SECRET"),
		GraphCloud:        envOr("GRAPH_CLOUD", "commercial"),
		GraphMasterCalID:  os.Getenv("GRAPH_MASTER_CALENDAR_ID"),

		// Weather
		WeatherLat: envOr("WEATHER_LAT", "38.52"),
		WeatherLon: envOr("WEATHER_LON", "-77.29"),
	}

	// --- Degraded-feature warnings ---

	hasOpenAI := cfg.OpenAIKey != ""
	hasAzure := cfg.AzureOpenAIEndpoint != "" && cfg.AzureOpenAIKey != ""
	if !hasOpenAI && !hasAzure {
		slog.Warn("No AI provider configured — chat will use mock responses. Set OPENAI_API_KEY or AZURE_OPENAI_ENDPOINT + AZURE_OPENAI_KEY.")
	}

	hasGraph := cfg.GraphTenantID != "" && cfg.GraphClientID != "" && cfg.GraphClientSecret != ""
	if !hasGraph {
		slog.Warn("Microsoft Graph not configured — calendar, mail, SharePoint, and Teams will use mock data. Set GRAPH_TENANT_ID, GRAPH_CLIENT_ID, GRAPH_CLIENT_SECRET.")
	}

	return cfg
}

// envOr returns the value of the environment variable named key, or
// fallback if the variable is unset or empty.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
