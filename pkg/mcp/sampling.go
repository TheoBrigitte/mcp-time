package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// EnableSampling enables sampling capabilities for the server.
// This allows tools to request LLM completions from the client when needed.
//
// Note: Sampling is already enabled via server.WithSamplingCapabilities(true) in NewServer.
// This function is provided for documentation and potential custom sampling logic.
func EnableSampling(s *server.MCPServer) {
	s.EnableSampling()
}

// RequestTimeInterpretation is a helper function that uses sampling to interpret
// ambiguous time expressions by asking the LLM for clarification.
//
// This demonstrates how MCP servers can use sampling to enhance their capabilities
// by delegating complex interpretation tasks to the client's LLM.
func RequestTimeInterpretation(ctx context.Context, s *server.MCPServer, timeExpression string) (string, error) {
	request := mcp.CreateMessageRequest{
		CreateMessageParams: mcp.CreateMessageParams{
			Messages: []mcp.SamplingMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf(`You are helping interpret a time expression.

Time expression: "%s"

Please provide:
1. The most likely interpretation of this time expression
2. Any ambiguities that exist
3. The recommended canonical format (ISO 8601 / RFC3339)

Respond in a structured format:
Interpretation: <your interpretation>
Ambiguities: <any ambiguities>
Canonical format: <RFC3339 formatted time>
`, timeExpression),
					},
				},
			},
			MaxTokens: 500,
		},
	}

	result, err := s.RequestSampling(ctx, request)
	if err != nil {
		return "", fmt.Errorf("sampling request failed: %w", err)
	}

	// Extract text content from the result
	if textContent, ok := mcp.AsTextContent(result.Content); ok {
		return textContent.Text, nil
	}

	return "", fmt.Errorf("unexpected response content type")
}

// SamplingExample demonstrates how sampling can be used within tool handlers.
// This is an example that shows the pattern for using sampling in your tools.
//
// Note: This is not registered as a tool by default. It's here for reference.
func SamplingExample(s *server.MCPServer) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		inputText := request.GetString("text", "")

		// Use sampling to get LLM interpretation
		interpretation, err := RequestTimeInterpretation(ctx, s, inputText)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to interpret time: %v", err)), nil
		}

		return mcp.NewToolResultText(interpretation), nil
	}
}
