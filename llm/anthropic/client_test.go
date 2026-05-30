package anthropic

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

type Config struct {
	AnthropicAPIKey string
}

func LoadConfig() *Config {
	godotenv.Load("../../.env")
	return &Config{
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
	}
}

func TestNewNoMaxTokens(t *testing.T) {
	model := "claude-haiku-4-5"
	v := defaultMaxTokens

	client, err := NewAnthropicClient(
		LoadConfig().AnthropicAPIKey,
		model,
		nil,
	)
	if err != nil || client == nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if *client.MaxTokens != v {
		t.Fatalf("failed to create client with proper max tokens, got %d", client.MaxTokens)
	}
}

func TestRun(t *testing.T) {
	model := "claude-haiku-4-5"
	content := "Write a simple haiku about anything"

	client, err := NewAnthropicClient(
		LoadConfig().AnthropicAPIKey,
		model,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	response, err := client.Run(content)
	if err != nil {
		t.Fatalf("failed to execute inference: %v", err)
	}

	// if len(response) == 0 {
	// 	t.Fatalf("expected messages but got empty response")
	// }

	t.Logf("got response: %v", response)
}
