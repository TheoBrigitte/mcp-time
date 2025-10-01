package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterPrompts registers time-related prompts with the MCP server.
func RegisterPrompts(s *server.MCPServer) {
	// Time format helper prompt
	formatHelperPrompt := mcp.NewPrompt("time_format_helper",
		mcp.WithPromptDescription("Get help with understanding and using time format strings"),
		mcp.WithArgument("format_type",
			mcp.ArgumentDescription("Type of format: 'predefined' for list of built-in formats, 'custom' for custom layout guide"),
			mcp.RequiredArgument(),
		),
	)
	s.AddPrompt(formatHelperPrompt, TimeFormatHelper)

	// Timezone conversion prompt
	timezonePrompt := mcp.NewPrompt("timezone_conversion_guide",
		mcp.WithPromptDescription("Get guidance on converting times between timezones"),
		mcp.WithArgument("from_timezone",
			mcp.ArgumentDescription("Source timezone in IANA format (e.g., 'America/New_York')"),
		),
		mcp.WithArgument("to_timezone",
			mcp.ArgumentDescription("Target timezone in IANA format (e.g., 'Europe/London')"),
		),
	)
	s.AddPrompt(timezonePrompt, TimezoneConversionGuide)

	// Relative time expression helper
	relativeTimePrompt := mcp.NewPrompt("relative_time_examples",
		mcp.WithPromptDescription("Get examples of natural language relative time expressions"),
	)
	s.AddPrompt(relativeTimePrompt, RelativeTimeExamples)
}

// TimeFormatHelper provides information about time format strings.
func TimeFormatHelper(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	formatType := request.Params.Arguments["format_type"]

	var content string
	if formatType == "predefined" {
		content = fmt.Sprintf("## Predefined Time Formats\n\nThe following predefined formats are available:\n\n%s\n\nUsage example:\n- Use 'RFC3339' for: %s\n- Use 'Kitchen' for: %s\n- Use 'UnixDate' for: %s",
			formatDescription,
			"2025-10-01T15:04:05Z07:00",
			"3:04PM",
			"Mon Oct 1 15:04:05 UTC 2025",
		)
	} else if formatType == "custom" {
		content = fmt.Sprintf("## Custom Time Format Layout\n\n%s\n\n### Examples\n\n- '2006-01-02' produces: 2025-10-01\n- 'Jan 2, 2006' produces: Oct 1, 2025\n- '3:04 PM' produces: 3:04 PM\n- 'Monday, January 2, 2006 at 3:04:05 PM MST' produces: Wednesday, October 1, 2025 at 3:04:05 PM UTC",
			formatDescription,
		)
	} else {
		content = "## Time Format Help\n\nPlease specify a format_type:\n- 'predefined': List all predefined time formats\n- 'custom': Learn how to create custom time format layouts"
	}

	return &mcp.GetPromptResult{
		Description: "Time format help",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: content,
				},
			},
		},
	}, nil
}

// TimezoneConversionGuide provides guidance on timezone conversions.
func TimezoneConversionGuide(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	fromTz := request.Params.Arguments["from_timezone"]
	toTz := request.Params.Arguments["to_timezone"]

	var content string
	if fromTz != "" && toTz != "" {
		content = fmt.Sprintf(`## Timezone Conversion: %s → %s

To convert a time from %s to %s, use the 'convert_timezone' tool:

**Example:**
- Input: "2025-10-01 15:04:05"
- Input Timezone: "%s"
- Output Timezone: "%s"

**Common Use Cases:**
- Meeting scheduling across regions
- Coordinating events in different time zones
- Understanding time differences for international communication

**Important Notes:**
- Daylight Saving Time (DST) is automatically handled
- IANA timezone names are case-sensitive
- Use 'UTC' for Coordinated Universal Time
`, fromTz, toTz, fromTz, toTz, fromTz, toTz)
	} else {
		content = `## Timezone Conversion Guide

Timezones use the IANA Time Zone Database format (e.g., 'America/New_York', 'Europe/London', 'Asia/Tokyo').

**Popular Timezones:**
- America/New_York (EST/EDT)
- America/Los_Angeles (PST/PDT)
- America/Chicago (CST/CDT)
- Europe/London (GMT/BST)
- Europe/Paris (CET/CEST)
- Asia/Tokyo (JST)
- Australia/Sydney (AEDT/AEST)
- UTC (Coordinated Universal Time)

**Tips:**
- Use the 'convert_timezone' tool for conversions
- Specify both input and output timezones for clarity
- DST transitions are handled automatically
`
	}

	return &mcp.GetPromptResult{
		Description: "Timezone conversion guidance",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: content,
				},
			},
		},
	}, nil
}

// RelativeTimeExamples provides examples of natural language time expressions.
func RelativeTimeExamples(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	content := fmt.Sprintf(`## Relative Time Expression Examples

The 'relative_time' tool understands natural language expressions. Here are examples:

**Immediate:**
- "now" - Current time
- "today" - Start of today

**Past:**
- "yesterday" - Yesterday at the same time
- "5 minutes ago" - 5 minutes before now
- "three days ago" - 3 days in the past
- "last week" - Previous week
- "last month" - Previous month
- "last year" - Previous year

**Future:**
- "tomorrow" - Tomorrow at the same time
- "in 2 hours" - 2 hours from now
- "next week" - Following week
- "next month" - Following month
- "one year from now" - One year in the future

**Specific Times:**
- "yesterday at 10am" - Yesterday at 10:00 AM
- "last sunday at 5:30pm" - Previous Sunday at 5:30 PM
- "sunday at 22:45" - Coming/past Sunday at 22:45
- "next January" - Start of next January
- "December 25th at 7:30am" - Dec 25 at 7:30 AM

**Time Only (relative to reference time's date):**
- "10am" - 10:00 AM on the reference date
- "10:05pm" - 10:05 PM on the reference date
- "10:05:22pm" - 10:05:22 PM on the reference date

%s
`, relativeTimeDescription)

	return &mcp.GetPromptResult{
		Description: "Relative time expression examples",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: content,
				},
			},
		},
	}, nil
}
