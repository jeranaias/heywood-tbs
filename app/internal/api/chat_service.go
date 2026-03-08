package api

import (
	"log/slog"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// eastern is the US Eastern timezone for Quantico, VA.
var eastern *time.Location

func init() {
	var err error
	eastern, err = time.LoadLocation("America/New_York")
	if err != nil {
		eastern = time.FixedZone("EST", -5*3600)
	}
}

// nowET returns the current time in Eastern timezone.
func nowET() time.Time { return time.Now().In(eastern) }

// ChatService manages OpenAI/Azure OpenAI API interactions.
type ChatService struct {
	client  *openai.Client
	model   string // "gpt-4o" or Azure deployment name
	isAzure bool
}

// NewChatService creates a chat service with automatic Azure/OpenAI detection.
//
// Env vars checked:
//   - AZURE_OPENAI_ENDPOINT + OPENAI_API_KEY → Azure OpenAI
//   - OPENAI_API_KEY alone → Public OpenAI
//   - Neither → nil (mock mode)
func NewChatService() *ChatService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	azureEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	azureDeployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT")

	if azureEndpoint != "" && apiKey != "" {
		// Azure OpenAI (IL5-ready path)
		config := openai.DefaultAzureConfig(apiKey, azureEndpoint)
		config.APIVersion = "2024-12-01-preview"
		model := "gpt-4o"
		if azureDeployment != "" {
			model = azureDeployment
		}
		slog.Info("Azure OpenAI configured", "endpoint", azureEndpoint, "deployment", model)
		return &ChatService{
			client:  openai.NewClientWithConfig(config),
			model:   model,
			isAzure: true,
		}
	}

	if apiKey != "" {
		// Public OpenAI
		slog.Info("OpenAI configured (public API)")
		return &ChatService{
			client: openai.NewClient(apiKey),
			model:  string(openai.GPT4o),
		}
	}

	return nil // mock mode
}
