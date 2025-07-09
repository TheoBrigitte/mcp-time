package mcp

import (
	"context"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/TheoBrigitte/mcp-time/pkg/datetime"
)

// CurrentTime is the handler for the 'current_time' MCP tool.
// It returns the current time, optionally formatted and in a specific timezone.
func CurrentTime(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	timezone := request.GetString("timezone", "")
	format := request.GetString("format", "")

	output, err := datetime.CurrentTime(timezone, format)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

// ConvertTime is the handler for the 'convert_timezone' MCP tool.
// It converts a time from one timezone to another.
func ConvertTime(ctx context.Context, request mcp.CallToolRequest) (r *mcp.CallToolResult, err error) {
	inputTime := request.GetString("time", "")
	inputTimezone := request.GetString("input_timezone", "")
	outputTimezone := request.GetString("output_timezone", "")
	format := request.GetString("format", "")

	output, err := datetime.ConvertTime(inputTime, inputTimezone, outputTimezone, format)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

// TimeAdd is the handler for the 'add_time' MCP tool.
// It adds a duration to a given time.
func TimeAdd(ctx context.Context, request mcp.CallToolRequest) (r *mcp.CallToolResult, err error) {
	inputTime := request.GetString("time", "")
	duration := request.GetString("duration", "")
	timezone := request.GetString("timezone", "")
	format := request.GetString("format", "")

	output, err := datetime.TimeAdd(inputTime, duration, timezone, format)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

// RelativeTime is the handler for the 'relative_time' MCP tool.
// It parses a natural language time expression relative to a given time.
func RelativeTime(ctx context.Context, request mcp.CallToolRequest) (r *mcp.CallToolResult, err error) {
	inputTime := request.GetString("time", "")
	relativeTime := request.GetString("text", "")
	timezone := request.GetString("timezone", "")
	format := request.GetString("format", "")

	output, err := datetime.RelativeTime(inputTime, relativeTime, timezone, format)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(output), nil
}

// CompareTime is the handler for the 'compare_time' MCP tool.
// It compares two times and returns -1, 0, or 1.
func CompareTime(ctx context.Context, request mcp.CallToolRequest) (r *mcp.CallToolResult, err error) {
	timeA := request.GetString("time_a", "")
	timeB := request.GetString("time_b", "")

	result, err := datetime.CompareTime(timeA, timeB)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	output := strconv.Itoa(result)

	return mcp.NewToolResultText(output), nil
}
