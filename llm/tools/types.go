package tools

import "encoding/json"

type ToolDefinition struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	InputSchema   InputSchema         `json:"input_schema"`
	InputExamples []map[string]string `json:"input_examples,omitempty"`
}

type InputSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]ToolProperties `json:"properties"`
	Required   []string                  `json:"required,omitempty"`
}

type ToolProperties struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type ToolHandler struct {
	Definition  ToolDefinition
	FuncHandler func(raw json.RawMessage) (any, error)
}
