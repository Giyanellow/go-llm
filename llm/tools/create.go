package tools

import (
	"encoding/json"
	"fmt"
)

// scratch pad for optional fields - to be used for With... funcs
type toolConfig struct {
	inputExamples []map[string]string
}

// type wrapper for all
type ToolOption func(*toolConfig)

func WithInputExamples(examples ...map[string]string) ToolOption {
	return func(cfg *toolConfig) {
		cfg.inputExamples = append(cfg.inputExamples, examples...)
	}
}

// Sample use:
// addTool := CreateTool("addTool", "a tool that adds two numbers", schema, addFunc, WithInputExamples{...})
func CreateTool[T any, R any](name string, description string, schema InputSchema, handler func(T) (R, error), opts ...ToolOption) (ToolHandler, error) {
	if schema.Type != "object" {
		return ToolHandler{}, fmt.Errorf("schema type must be 'object'")
	}

	cfg := &toolConfig{}
	// loop over each opt func, append to cfg within opt func
	for _, opt := range opts {
		opt(cfg)
	}

	defn := ToolDefinition{
		Name:          name,
		Description:   description,
		InputSchema:   schema,
		InputExamples: cfg.inputExamples,
	}

	// FuncHandler is a function that does:
	// - takes in a raw json for any input of the assigned function
	// - converts json into proper struct provided by user
	// - returns function that has args assignment based from user given arg struct
	return ToolHandler{
		Definition: defn,
		FuncHandler: func(raw json.RawMessage) (any, error) {
			var args T
			err := json.Unmarshal(raw, &args)
			if err != nil {
				return nil, err
			}
			return handler(args)
		},
	}, nil
}
