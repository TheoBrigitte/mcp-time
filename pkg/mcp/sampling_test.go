package mcp

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func TestEnableSampling(t *testing.T) {
	// Test that EnableSampling can be called without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("EnableSampling() panicked: %v", r)
		}
	}()

	s := server.NewMCPServer("test-server", "1.0.0")
	EnableSampling(s)
}

func TestRequestTimeInterpretation(t *testing.T) {
	// This test verifies the function structure but cannot fully test
	// the sampling functionality without a connected client.
	// We test that the function can be called and returns appropriate errors.

	ctx := context.Background()
	s := server.NewMCPServer("test-server", "1.0.0")
	EnableSampling(s)

	// This should fail because there's no client connected to handle sampling
	_, err := RequestTimeInterpretation(ctx, s, "yesterday at 3pm")

	// We expect an error since there's no client to handle the sampling request
	if err == nil {
		t.Log("RequestTimeInterpretation() did not error (client may be connected)")
	} else {
		// This is the expected case - no client connected
		t.Logf("RequestTimeInterpretation() returned expected error: %v", err)
	}
}

func TestSamplingExample(t *testing.T) {
	// Test that SamplingExample returns a valid handler function
	s := server.NewMCPServer("test-server", "1.0.0")
	EnableSampling(s)

	handler := SamplingExample(s)
	if handler == nil {
		t.Fatal("SamplingExample() returned nil handler")
	}

	// Test calling the handler
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "test_tool",
			Arguments: map[string]interface{}{
				"text": "tomorrow",
			},
		},
	}

	// This should fail because there's no client connected
	result, err := handler(ctx, request)

	// We expect either an error or a tool result with an error message
	if err != nil {
		t.Logf("Handler returned error (expected without client): %v", err)
	} else if result != nil {
		t.Logf("Handler returned result: %+v", result)
	}
}

func TestSamplingCreateMessageRequest(t *testing.T) {
	// Test that we can construct a valid CreateMessageRequest
	request := mcp.CreateMessageRequest{
		CreateMessageParams: mcp.CreateMessageParams{
			Messages: []mcp.SamplingMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: "Test message",
					},
				},
			},
			MaxTokens: 100,
		},
	}

	if len(request.CreateMessageParams.Messages) != 1 {
		t.Errorf("CreateMessageRequest has %d messages, want 1", len(request.CreateMessageParams.Messages))
	}

	if request.CreateMessageParams.MaxTokens != 100 {
		t.Errorf("CreateMessageRequest MaxTokens = %d, want 100", request.CreateMessageParams.MaxTokens)
	}

	message := request.CreateMessageParams.Messages[0]
	if message.Role != mcp.RoleUser {
		t.Errorf("Message role = %v, want %v", message.Role, mcp.RoleUser)
	}

	textContent, ok := mcp.AsTextContent(message.Content)
	if !ok {
		t.Fatal("Message content is not text content")
	}

	if textContent.Text != "Test message" {
		t.Errorf("Message text = %q, want %q", textContent.Text, "Test message")
	}
}

func TestSamplingIntegration(t *testing.T) {
	// Integration test to verify sampling is enabled in the server
	s := NewServer("test-server", "1.0.0")
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// Verify that sampling was enabled during server initialization
	// The server should have sampling capability enabled
	// Note: We can't directly test this without inspecting server internals,
	// but we can verify the server was created successfully
}
