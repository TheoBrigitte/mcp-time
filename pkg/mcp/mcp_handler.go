// Package mcp provides MCP handlers for time and date utilities.
// It enables AI assistants to perform time-related operations through standardized MCP tools.
package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/TheoBrigitte/mcp-time/pkg/datetime"
)

// formatDescription provides detailed information about the supported time formats for MCP tools.
var formatDescription = `The output time format, which can be a predefined format or a custom layout.

## Predefined Formats

` + fmt.Sprintf("`%s`", strings.Join(datetime.GetFormats(), "`, `")) + `

## Custom Format

A custom format can be built using the following components. Each component shows an example of how a part of the reference time is formatted. Only these values are recognized. Any text in the layout string that is not a recognized component will be treated as a literal.

- Year: "2006", "06"
- Month: "Jan", "January", "01", "1"
- Day of the week: "Mon", "Monday"
- Day of the month: "2", "_2", "02"
- Day of the year: "__2", "002"
- Hour: "15", "3", "03" (PM or AM)
- Minute: "4", "04"
- Second: "5", "05"
- AM/PM mark: "PM"

### Numeric Time Zone Offsets

- "-0700"     (±hhmm)
- "-07:00"    (±hh:mm)
- "-07"       (±hh)
- "-070000"   (±hhmmss)
- "-07:00:00" (±hh:mm:ss)

Replacing the sign with a "Z" triggers ISO 8601 behavior, which prints "Z" for the UTC zone:

- "Z0700"      (Z or ±hhmm)
- "Z07:00"     (Z or ±hh:mm)
- "Z07"        (Z or ±hh)
- "Z070000"    (Z or ±hhmmss)
- "Z07:00:00"  (Z or ±hh:mm:ss)

Within the format string, the underscores in "_2" and "__2" represent spaces that may be replaced by digits if the following number has multiple digits, for compatibility with fixed-width Unix time formats. A leading zero represents a zero-padded value.
The formats __2 and 002 are space-padded and zero-padded three-character day of year; there is no unpadded day of year format.

### Fractional Seconds

A comma or decimal point followed by one or more zeros represents a fractional second, printed to the given number of decimal places. A comma or decimal point followed by one or more nines represents a fractional second with trailing zeros removed.
For example, "15:04:05.000" formats or parses with millisecond precision.`

// durationDescription explains the format for duration strings used in MCP tools.
const durationDescription = `The duration to add or subtract. Use a negative value to subtract.
Examples:
- "1h2m3s" to add 1 hour, 2 minutes, and 3 seconds.
- "-1h" to subtract 1 hour.`

// relativeTimeDescription provides examples of natural language expressions for relative time.
const relativeTimeDescription = `A relative time expression in natural language.
Examples:
- "now"
- "today"
- "yesterday"
- "5 minutes ago"
- "three days ago"
- "last month"
- "next month"
- "one year from now"
- "yesterday at 10am"
- "last sunday at 5:30pm"
- "sunday at 22:45"
- "next January"
- "last February"
- "December 25th at 7:30am"
- "10am"
- "10:05pm"
- "10:05:22pm"`

// compareDescription explains the output of the compare_time tool.
const compareDescription = `Compares two times. Returns -1 if the first time is before the second, 0 if they are equal, and 1 if the first time is after the second.`

// RegisterHandlers registers the time and date MCP tools, prompts, and resources with the provided MCP server.
//
// Parameters:
//   - s: The MCP server instance to register tools, prompts, and resources with.
//
// Returns an error if tool registration fails (though current implementation always returns nil).
func RegisterHandlers(s *server.MCPServer) {
	// Register prompts
	RegisterPrompts(s)

	// Register resources
	RegisterResources(s)

	// Enable sampling (already configured in server capabilities, but explicitly enable)
	EnableSampling(s)

	// Register tools
	currentTime := mcp.NewTool("current_time",
		mcp.WithDescription("Returns the current time."),
		mcp.WithString("format",
			mcp.Description(formatDescription),
			mcp.DefaultString(datetime.GetDefaultFormat()),
		),
		timezoneProperty,

		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	)
	s.AddTool(currentTime, CurrentTime)

	convertTimezone := mcp.NewTool("convert_timezone",
		mcp.WithDescription("Converts a time from one timezone to another."),
		mcp.WithString("input_timezone",
			mcp.Description("The timezone of the input time, in IANA format (e.g., 'America/New_York'). If the input time string contains a timezone, it will take precedence."),
			mcp.DefaultString(datetime.GetDefaultTimezone()),
		),
		mcp.WithString("output_timezone",
			mcp.Description("The target timezone for the output, in IANA format (e.g., 'America/New_York')."),
			mcp.DefaultString(datetime.GetDefaultTimezone()),
		),
		timeProperty,
		formatProperty,

		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	)
	s.AddTool(convertTimezone, ConvertTime)

	addTime := mcp.NewTool("add_time",
		mcp.WithDescription("Adds or subtracts a duration from a given time."),
		mcp.WithString("duration",
			mcp.Description(durationDescription),
			mcp.Required(),
		),
		timeProperty,
		timezoneProperty,
		formatProperty,

		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	)
	s.AddTool(addTime, TimeAdd)

	relativeTime := mcp.NewTool("relative_time",
		mcp.WithDescription("Returns a time based on a relative natural language expression."),
		mcp.WithString("text",
			mcp.Description(relativeTimeDescription),
			mcp.Required(),
		),
		timeProperty,
		timezoneProperty,
		formatProperty,

		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	)
	s.AddTool(relativeTime, RelativeTime)

	compareTime := mcp.NewTool("compare_time",
		mcp.WithDescription(compareDescription),
		mcp.WithString("time_a",
			mcp.Description("The first time to compare."),
			mcp.Required(),
		),
		mcp.WithString("time_a_timezone",
			mcp.Description("Timezone form time_a, in IANA format (e.g., 'America/New_York')."),
			mcp.DefaultString(datetime.GetDefaultTimezone()),
		),
		mcp.WithString("time_b",
			mcp.Description("The second time to compare."),
			mcp.Required(),
		),
		mcp.WithString("time_b_timezone",
			mcp.Description("Timezone for time_b, in IANA format (e.g., 'America/New_York')."),
			mcp.DefaultString(datetime.GetDefaultTimezone()),
		),

		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithIdempotentHintAnnotation(true),
		mcp.WithOpenWorldHintAnnotation(false),
	)
	s.AddTool(compareTime, CompareTime)
}
