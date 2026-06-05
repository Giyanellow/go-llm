package anthropic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giyanellow/go-llm/internal/httputil"
	"github.com/giyanellow/go-llm/llm"
	"github.com/giyanellow/go-llm/llm/tools"
)

const defaultMaxTokens = 1024

type clientConfig struct {
	tools []tools.ToolDefinition
}

type ClientOption func(*clientConfig)

func WithTools(tools ...tools.ToolDefinition) ClientOption {
	return func(cfg *clientConfig) {
		cfg.tools = append(cfg.tools, tools...)
	}
}

type AnthropicClient struct {
	apiKey     string
	Model      string
	Tools      []tools.ToolDefinition
	MaxTokens  *int
	httpClient *http.Client
}

func NewAnthropicClient(apiKey string, Model string, MaxTokens *int, opts ...ClientOption) (*AnthropicClient, error) {
	var err error
	if apiKey == "" {
		err = fmt.Errorf("API Key must be satisfied")
	} else if Model == "" {
		err = fmt.Errorf("Model must be set")
	}

	if MaxTokens == nil {
		v := defaultMaxTokens
		MaxTokens = &v
	}

	// if expected to return pointer, can return nil for errors
	if err != nil {
		return nil, err
	}

	cfg := &clientConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	return &AnthropicClient{
		apiKey:     apiKey,
		Model:      Model,
		Tools:      cfg.tools,
		MaxTokens:  MaxTokens,
		httpClient: &http.Client{},
	}, nil
}

func (c *AnthropicClient) Run(prompt string, specificTools ...tools.ToolDefinition) ([]llm.Message, error) {
	// create input struct
	messageInput := anthropicMessage{
		Role:    "user",
		Content: prompt,
	}

	messagesPayload := []anthropicMessage{messageInput}

	requestBody := anthropicRequestBody{
		Model:     c.Model,
		Input:     messagesPayload,
		MaxTokens: *c.MaxTokens,
		Tools:     []tools.ToolDefinition{},
	}

	// fill blank anthropicRequestBody.tools
	for _, tool := range specificTools {
		requestBody.Tools = append(requestBody.Tools, tool)
	}

	requestHeaders := map[string]string{
		"content-type":      "application/json",
		"x-api-key":         c.apiKey,
		"anthropic-version": "2023-06-01",
	}

	resp, err := httputil.Do(c.httpClient, chatBaseUrl, "POST", requestHeaders, requestBody)
	if err != nil {
		return []llm.Message{}, err
	}

	var llmResponse anthropicResponseBody
	err = json.Unmarshal(resp, &llmResponse)
	if err != nil {
		return []llm.Message{}, fmt.Errorf("Error in parsing response: %s", err.Error())
	}

	var messages []llm.Message

	for _, item := range llmResponse.Content {
		messages = append(messages, llm.Message{
			Role:    llmResponse.Role,
			Content: item.Text,
			Type:    item.MessageType,
		})
	}

	return messages, nil
}

func (c *AnthropicClient) GetModel() string {
	return c.Model
}
