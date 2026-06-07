package anthropic

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/giyanellow/go-llm/internal/httputil"
	"github.com/giyanellow/go-llm/llm"
	"github.com/giyanellow/go-llm/llm/tools"
)

const defaultMaxTokens = 1024

type clientConfig struct {
	toolDefinitions []tools.ToolDefinition
	toolHandlers    map[string]func(json.RawMessage) (any, error)
}

type ClientOption func(*clientConfig)

func WithClientTools(tools ...tools.ToolHandler) ClientOption {
	return func(cfg *clientConfig) {
		if cfg.toolHandlers == nil {
			cfg.toolHandlers = make(map[string]func(json.RawMessage) (any, error))
		}
		for _, handler := range tools {
			cfg.toolDefinitions = append(cfg.toolDefinitions, handler.Definition)
			cfg.toolHandlers[handler.Definition.Name] = handler.FuncHandler
		}
	}
}

type AnthropicClient struct {
	apiKey          string
	Model           string
	SystemPrompt    string
	ToolDefinitions []tools.ToolDefinition
	toolHandlers    map[string]func(json.RawMessage) (any, error)
	MaxTokens       *int
	httpClient      *http.Client
	clientHeaders   map[string]string
}

func NewAnthropicClient(apiKey string, Model string, SystemPrompt string, MaxTokens *int, opts ...ClientOption) (*AnthropicClient, error) {
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

	requestHeaders := map[string]string{
		"content-type":      "application/json",
		"x-api-key":         apiKey,
		"anthropic-version": "2023-06-01",
	}

	return &AnthropicClient{
		apiKey:          apiKey,
		Model:           Model,
		SystemPrompt:    SystemPrompt,
		ToolDefinitions: cfg.toolDefinitions,
		toolHandlers:    cfg.toolHandlers,
		MaxTokens:       MaxTokens,
		httpClient:      &http.Client{},
		clientHeaders:   requestHeaders,
	}, nil
}

func (c *AnthropicClient) Run(prompt string, specificTools ...tools.ToolHandler) ([]llm.Message, error) {
	// create input struct
	content, err := json.Marshal(prompt)

	messageInput := anthropicMessage{
		Role:    "user",
		Content: json.RawMessage(content),
	}

	messagesPayload := []anthropicMessage{messageInput}

	requestBody := anthropicRequestBody{
		Model:           c.Model,
		SystemPrompt:    c.SystemPrompt,
		Input:           messagesPayload,
		MaxTokens:       *c.MaxTokens,
		ToolsDefinition: c.ToolDefinitions,
	}

	// fill blank anthropicRequestBody.tools
	for _, tool := range specificTools {
		requestBody.ToolsDefinition = append(requestBody.ToolsDefinition, tool.Definition)
	}

	resp, err := httputil.Do(c.httpClient, chatBaseUrl, "POST", c.clientHeaders, requestBody)
	if err != nil {
		return []llm.Message{}, err
	}

	var llmResponse anthropicResponseBody
	err = json.Unmarshal(resp, &llmResponse)
	if err != nil {
		return []llm.Message{}, fmt.Errorf("Error in parsing response: %s", err.Error())
	}

	var messages []llm.Message

	log.Printf("api call response: %v", string(resp))

	lastMessage := llmResponse.Content[len(llmResponse.Content)-1]

	// loop for tool calling
	for llmResponse.StopReason == "tool_use" {
		assistantContent, err := json.Marshal(llmResponse.Content)
		messagesPayload = append(messagesPayload, anthropicMessage{
			Role:    "assistant",
			Content: json.RawMessage(assistantContent),
		})

		toolResult, err := c.handleToolUse(lastMessage.Name, lastMessage.ID, lastMessage.Input)
		// custom check from anthropicToolResultMessage for err
		if err != nil {
			log.Printf("there was an error running the tool: %v", err)
		}

		toolResultContent, err := json.Marshal([]anthropicToolResultMessage{toolResult})
		if err != nil {
			return []llm.Message{}, err
		}

		messagesPayload = append(messagesPayload, anthropicMessage{
			Role:    "user",
			Content: json.RawMessage(toolResultContent),
		})

		requestWithToolResult := requestBody
		requestWithToolResult.Input = messagesPayload

		// API call with tool result
		resp, err := httputil.Do(c.httpClient, chatBaseUrl, "POST", c.clientHeaders, requestWithToolResult)
		if err != nil {
			return []llm.Message{}, err
		}

		err = json.Unmarshal(resp, &llmResponse)
		if err != nil {
			return []llm.Message{}, fmt.Errorf("Error in parsing response: %s", err.Error())
		}

		lastMessage = llmResponse.Content[len(llmResponse.Content)-1]

	}

	for _, item := range llmResponse.Content {
		messages = append(messages, llm.Message{
			Role:    llmResponse.Role,
			Content: item.Text,
			Type:    item.MessageType,
		})
	}

	return messages, nil
}

func (c *AnthropicClient) handleToolUse(toolName string, toolID string, toolInput json.RawMessage) (anthropicToolResultMessage, error) {
	funcToUse := c.toolHandlers[toolName]
	res, err := funcToUse(toolInput)
	if err != nil {
		return anthropicToolResultMessage{
			MessageType: "tool_result",
			ToolUseID:   toolID,
			Content:     fmt.Sprintf("Tool Call Error: %v", err.Error()),
			IsError:     true,
		}, fmt.Errorf("there was an error running your tool: %v", err)
	}
	safeResponse := fmt.Sprintf("Here is the tool result: %v", res)

	toolMessage := anthropicToolResultMessage{
		MessageType: "tool_result",
		ToolUseID:   toolID,
		Content:     safeResponse,
	}

	return toolMessage, nil
}

func (c *AnthropicClient) GetModel() string {
	return c.Model
}
