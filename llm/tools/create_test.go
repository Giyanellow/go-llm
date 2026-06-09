package tools

import (
	"reflect"
	"testing"
)

type sampleToolArgs struct {
	A string `json:"a"`
	B string `json:"b"`
}

var testSchema = InputSchema{
	Type: "object",
	Properties: map[string]ToolProperties{
		"a": {
			Type:        "string",
			Description: "any sample string",
		},
		"b": {
			Type:        "int",
			Description: "any sample int",
		},
	},
}

var sampleTool = func(args sampleToolArgs) (string, error) {
	return "sample return", nil
}

var invalidTestSchema = InputSchema{
	Type: "string",
	Properties: map[string]ToolProperties{
		"a": {
			Type:        "string",
			Description: "any sample string",
		},
		"b": {
			Type:        "int",
			Description: "any sample int",
		},
	},
}

var (
	toolName        = "test-tool"
	toolDescription = "a tool for testing"
)

func TestCreateTool(t *testing.T) {
	testTool, err := CreateTool(toolName, toolDescription, testSchema, sampleTool)
	if err != nil {
		t.Fatalf("failed to create test tool: %v", err)
	}
	if testTool.Definition.Name != toolName {
		t.Fatalf("test tool name is not equal to tool name")
	}
	if testTool.Definition.Description != toolDescription {
		t.Fatalf("test tool name is not equal to tool description")
	}
	if !reflect.DeepEqual(testTool.Definition.InputSchema, testSchema) {
		t.Fatalf("test tool input schema not equal to test schema")
	}
	if testTool.FuncHandler == nil {
		t.Fatalf("test tool handler was set, got nil")
	}
}

func TestNotTypeObjectTool(t *testing.T) {
	_, err := CreateTool(toolName, toolDescription, invalidTestSchema, sampleTool)
	if err == nil {
		t.Fatalf("expected error for invalid schema, got nil")
	}
}
