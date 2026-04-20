package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/TheoBrigitte/mcp-time/pkg/datetime"
)

// RegisterResources registers time-related resources with the MCP server.
func RegisterResources(s *server.MCPServer) {
	// Available time formats resource
	formatsResource := mcp.NewResource(
		"time://formats",
		"Time Formats",
		mcp.WithResourceDescription("List of all available predefined time formats"),
		mcp.WithMIMEType("text/plain"),
	)
	s.AddResource(formatsResource, TimeFormatsResource)

	// Current time in various formats - template resource
	currentTimeTemplate := mcp.NewResourceTemplate(
		"time://current/{timezone}",
		"Current Time",
		mcp.WithTemplateDescription("Current time in the specified timezone. Use IANA timezone format (e.g., 'America/New_York', 'Europe/London', 'UTC')"),
		mcp.WithTemplateMIMEType("text/plain"),
	)
	s.AddResourceTemplate(currentTimeTemplate, CurrentTimeResource)

	// Timezone information resource template
	timezoneInfoTemplate := mcp.NewResourceTemplate(
		"time://timezone-info/{timezone}",
		"Timezone Information",
		mcp.WithTemplateDescription("Information about a specific timezone including current time and UTC offset. Use IANA timezone format."),
		mcp.WithTemplateMIMEType("text/markdown"),
	)
	s.AddResourceTemplate(timezoneInfoTemplate, TimezoneInfoResource)

	// Popular timezones list
	popularTimezonesResource := mcp.NewResource(
		"time://timezones/popular",
		"Popular Timezones",
		mcp.WithResourceDescription("List of commonly used IANA timezone names"),
		mcp.WithMIMEType("text/markdown"),
	)
	s.AddResource(popularTimezonesResource, PopularTimezonesResource)

	// Relative time expressions guide
	relativeTimeGuideResource := mcp.NewResource(
		"time://guide/relative-expressions",
		"Relative Time Expressions Guide",
		mcp.WithResourceDescription("Guide to natural language relative time expressions"),
		mcp.WithMIMEType("text/markdown"),
	)
	s.AddResource(relativeTimeGuideResource, RelativeTimeGuideResource)
}

// TimeFormatsResource returns a list of all available predefined time formats.
func TimeFormatsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	formats := datetime.GetFormats()
	content := "Available predefined time formats:\n\n" + strings.Join(formats, "\n")

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     content,
		},
	}, nil
}

// CurrentTimeResource returns the current time in the specified timezone.
func CurrentTimeResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract timezone from URI path (time://current/{timezone})
	// The timezone may contain slashes (e.g., America/New_York)
	prefix := "time://current/"
	if !strings.HasPrefix(request.Params.URI, prefix) {
		return nil, fmt.Errorf("invalid URI format: expected time://current/{timezone}")
	}
	timezone := request.Params.URI[len(prefix):]
	if timezone == "" {
		return nil, fmt.Errorf("invalid URI format: timezone is required")
	}

	// Get current time in the specified timezone with default format
	currentTime, err := datetime.CurrentTime(timezone, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get current time: %w", err)
	}

	content := fmt.Sprintf("Current time in %s:\n%s", timezone, currentTime)

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     content,
		},
	}, nil
}

// TimezoneInfoResource provides detailed information about a timezone.
func TimezoneInfoResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Extract timezone from URI path (time://timezone-info/{timezone})
	// The timezone may contain slashes (e.g., Asia/Tokyo)
	prefix := "time://timezone-info/"
	if !strings.HasPrefix(request.Params.URI, prefix) {
		return nil, fmt.Errorf("invalid URI format: expected time://timezone-info/{timezone}")
	}
	timezone := request.Params.URI[len(prefix):]
	if timezone == "" {
		return nil, fmt.Errorf("invalid URI format: timezone is required")
	}

	// Get current time in various formats
	currentTime, err := datetime.CurrentTime(timezone, "RFC3339")
	if err != nil {
		return nil, fmt.Errorf("failed to get timezone info: %w", err)
	}

	utcTime, err := datetime.CurrentTime("UTC", "RFC3339")
	if err != nil {
		return nil, fmt.Errorf("failed to get UTC time: %w", err)
	}

	content := fmt.Sprintf(`# Timezone Information: %s

## Current Time
- Local: %s
- UTC: %s

## Timezone Details
- IANA Name: %s
- Format: IANA Time Zone Database format

## Usage
Use this timezone identifier with time conversion tools to work with times in this region.
`, timezone, currentTime, utcTime, timezone)

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     content,
		},
	}, nil
}

// PopularTimezonesResource returns a list of commonly used timezones.
func PopularTimezonesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	content := `# Popular Timezones

## North America
- **America/New_York** - Eastern Time (EST/EDT)
- **America/Chicago** - Central Time (CST/CDT)
- **America/Denver** - Mountain Time (MST/MDT)
- **America/Los_Angeles** - Pacific Time (PST/PDT)
- **America/Phoenix** - Mountain Time (no DST)
- **America/Anchorage** - Alaska Time (AKST/AKDT)
- **Pacific/Honolulu** - Hawaii Time (HST)

## Europe
- **Europe/London** - Greenwich Mean Time (GMT/BST)
- **Europe/Paris** - Central European Time (CET/CEST)
- **Europe/Berlin** - Central European Time (CET/CEST)
- **Europe/Rome** - Central European Time (CET/CEST)
- **Europe/Madrid** - Central European Time (CET/CEST)
- **Europe/Amsterdam** - Central European Time (CET/CEST)
- **Europe/Moscow** - Moscow Time (MSK)

## Asia
- **Asia/Dubai** - Gulf Standard Time (GST)
- **Asia/Kolkata** - India Standard Time (IST)
- **Asia/Shanghai** - China Standard Time (CST)
- **Asia/Tokyo** - Japan Standard Time (JST)
- **Asia/Hong_Kong** - Hong Kong Time (HKT)
- **Asia/Singapore** - Singapore Time (SGT)
- **Asia/Seoul** - Korea Standard Time (KST)

## Oceania
- **Australia/Sydney** - Australian Eastern Time (AEDT/AEST)
- **Australia/Melbourne** - Australian Eastern Time (AEDT/AEST)
- **Australia/Brisbane** - Australian Eastern Time (no DST)
- **Australia/Perth** - Australian Western Time (AWST)
- **Pacific/Auckland** - New Zealand Time (NZDT/NZST)

## Universal
- **UTC** - Coordinated Universal Time

## Notes
- Timezones with DST (Daylight Saving Time) automatically adjust
- Use IANA timezone names exactly as shown (case-sensitive)
- For more timezones, refer to the IANA Time Zone Database
`

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     content,
		},
	}, nil
}

// RelativeTimeGuideResource provides a guide to relative time expressions.
func RelativeTimeGuideResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	content := `# Relative Time Expressions Guide

Natural language time expressions that can be used with the relative_time tool.

## Current Time
- **now** - Current moment
- **today** - Start of today (00:00:00)

## Past Expressions
### Recent Past
- **yesterday** - Same time yesterday
- **X minutes ago** - Minutes in the past (e.g., "5 minutes ago", "30 minutes ago")
- **X hours ago** - Hours in the past (e.g., "2 hours ago", "12 hours ago")

### Days
- **X days ago** - Days in the past (e.g., "3 days ago", "seven days ago")
- **last week** - Previous week
- **last [day]** - Previous occurrence of day (e.g., "last Monday", "last Friday")

### Months and Years
- **last month** - Previous month
- **last year** - Previous year
- **last [month]** - Previous occurrence of month (e.g., "last January", "last December")

## Future Expressions
### Near Future
- **tomorrow** - Same time tomorrow
- **in X minutes** - Minutes in the future (e.g., "in 15 minutes")
- **in X hours** - Hours in the future (e.g., "in 2 hours")

### Days
- **in X days** - Days in the future (e.g., "in 3 days")
- **next week** - Following week
- **next [day]** - Next occurrence of day (e.g., "next Monday")

### Months and Years
- **next month** - Following month
- **next year** - Following year
- **one year from now** - Exactly one year in the future
- **next [month]** - Next occurrence of month (e.g., "next January")

## Specific Times
### With Time of Day
- **yesterday at 10am** - Yesterday at 10:00 AM
- **tomorrow at 3pm** - Tomorrow at 3:00 PM
- **last Sunday at 5:30pm** - Previous Sunday at 5:30 PM
- **next Friday at 9:00am** - Coming Friday at 9:00 AM

### Date-Specific
- **[Day] at [time]** - Specific day and time (e.g., "Sunday at 22:45")
- **[Month] [day] at [time]** - Specific date and time (e.g., "December 25th at 7:30am")

### Time Only
When you specify only a time, it's interpreted relative to the reference date:
- **10am** - 10:00 AM on the reference date
- **10:05pm** - 10:05 PM on the reference date
- **10:05:22pm** - 10:05:22 PM with seconds

## Tips
- You can use numbers ("5 days ago") or words ("five days ago")
- Times can be in 12-hour (with am/pm) or 24-hour format
- Expressions are case-insensitive
- The tool handles context automatically (e.g., "next Monday" from Friday vs from Wednesday)

## Examples in Context
` + "```" + `
"three days ago" -> 3 days before the reference time
"next January" -> The upcoming January from the reference time
"yesterday at 10am" -> 10:00 AM on the day before reference time
"in 2 hours" -> 2 hours after the reference time
` + "```" + `
`

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     content,
		},
	}, nil
}
