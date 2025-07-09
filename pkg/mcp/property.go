package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/TheoBrigitte/mcp-time/pkg/datetime"
)

var (
	// timeProperty is a reusable MCP property for a time string input.
	// It defaults to the current time if not provided.
	timeProperty = mcp.WithString("time",
		mcp.Description("Time in any format. Defaults to the current time."),
	)

	// formatProperty is a reusable MCP property for the output time format.
	formatProperty = mcp.WithString("format",
		mcp.Description("Output time format. See the 'current_time' tool for detailed format options."),
		mcp.DefaultString(datetime.GetDefaultFormat()),
	)

	// timezoneProperty is a reusable MCP property for specifying a timezone.
	timezoneProperty = mcp.WithString("timezone",
		mcp.Description("The target timezone for the output, in IANA format (e.g., 'America/New_York')."),
		mcp.DefaultString(datetime.GetDefaultTimezone()),
	)
)
