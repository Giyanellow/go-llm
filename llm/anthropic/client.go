package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/giyanellow/go-llm/llm"
)

// community practice to declared unused var to satisfy LSP
var _ llm.LLMClient = (*AnthropicClient)(nil)

const defaultMaxTokens = 1024

type AnthropicClient struct {
	apiKey     string
	Model      string
	MaxTokens  *int
	httpClient *http.Client
}

func NewAnthropicClient(apiKey string, Model string, MaxTokens *int) (*AnthropicClient, error) {
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

	return &AnthropicClient{
		apiKey:     apiKey,
		Model:      Model,
		MaxTokens:  MaxTokens,
		httpClient: &http.Client{},
	}, nil
}

func (c *AnthropicClient) Run(prompt string) ([]llm.Message, error) {
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
	}

	jsonRequestData, err := json.Marshal(requestBody)
	if err != nil {
		return []llm.Message{}, err
	}

	req, err := http.NewRequest(
		"POST",
		chatBaseUrl,
		bytes.NewBuffer(jsonRequestData),
	)
	if err != nil {
		return []llm.Message{}, err
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []llm.Message{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return []llm.Message{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []llm.Message{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var llmResponse anthropicResponseBody
	err = json.Unmarshal(body, &llmResponse)
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
