package tools

import "fmt"

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
// addTool := CreateTool("addTool", "a tool that adds two numbers", schema, WithInputExamples{...})
func CreateTool(name string, description string, schema InputSchema, opts ...ToolOption) (ToolDefinition, error) {
	if schema.Type != "object" {
		return ToolDefinition{}, fmt.Errorf("schema type must be 'object'")
	}

	cfg := &toolConfig{}
	// loop over each opt func, append to cfg within opt func
	for _, opt := range opts {
		opt(cfg)
	}
	return ToolDefinition{
		Name:          name,
		Description:   description,
		InputSchema:   schema,
		InputExamples: cfg.inputExamples,
	}, nil
}
