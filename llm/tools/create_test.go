package tools

import (
	"reflect"
	"testing"
)

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
	testTool, err := CreateTool(toolName, toolDescription, testSchema)
	if err != nil {
		t.Fatalf("failed to create test tool: %v", err)
	}
	if testTool.Name != toolName {
		t.Fatalf("test tool name is not equal to tool name")
	}
	if testTool.Description != toolDescription {
		t.Fatalf("test tool name is not equal to tool description")
	}
	if !reflect.DeepEqual(testTool.InputSchema, testSchema) {
		t.Fatalf("test tool input schema not equal to test schema")
	}
}

func TestNotTypeObjectTool(t *testing.T) {
	_, err := CreateTool(toolName, toolDescription, invalidTestSchema)
	if err == nil {
		t.Fatalf("expected error for invalid schema, got nil")
	}
}
