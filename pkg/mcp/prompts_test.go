package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestTimeFormatHelper(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		formatType     string
		expectContains []string
	}{
		{
			name:       "predefined formats",
			formatType: "predefined",
			expectContains: []string{
				"Predefined Time Formats",
				"RFC3339",
				"Kitchen",
				"UnixDate",
			},
		},
		{
			name:       "custom formats",
			formatType: "custom",
			expectContains: []string{
				"Custom Time Format Layout",
				"2006-01-02",
				"Jan 2, 2006",
				"3:04 PM",
			},
		},
		{
			name:       "no format type specified",
			formatType: "",
			expectContains: []string{
				"Time Format Help",
				"predefined",
				"custom",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.GetPromptRequest{
				Params: mcp.GetPromptParams{
					Name: "time_format_helper",
					Arguments: map[string]string{
						"format_type": tt.formatType,
					},
				},
			}

			result, err := TimeFormatHelper(ctx, request)
			if err != nil {
				t.Fatalf("TimeFormatHelper() error = %v", err)
			}

			if result == nil {
				t.Fatal("TimeFormatHelper() returned nil result")
			}

			if len(result.Messages) == 0 {
				t.Fatal("TimeFormatHelper() returned no messages")
			}

			// Check that the result contains expected content
			message := result.Messages[0]
			if textContent, ok := mcp.AsTextContent(message.Content); ok {
				content := textContent.Text
				for _, expected := range tt.expectContains {
					if !strings.Contains(content, expected) {
						t.Errorf("TimeFormatHelper() content missing expected string %q", expected)
					}
				}
			} else {
				t.Fatal("TimeFormatHelper() did not return text content")
			}
		})
	}
}

func TestTimezoneConversionGuide(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		fromTimezone   string
		toTimezone     string
		expectContains []string
	}{
		{
			name:         "with specific timezones",
			fromTimezone: "America/New_York",
			toTimezone:   "Europe/London",
			expectContains: []string{
				"Timezone Conversion",
				"America/New_York",
				"Europe/London",
				"convert_timezone",
			},
		},
		{
			name:         "without timezones",
			fromTimezone: "",
			toTimezone:   "",
			expectContains: []string{
				"Timezone Conversion Guide",
				"IANA Time Zone Database",
				"America/New_York",
				"Europe/London",
				"Asia/Tokyo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.GetPromptRequest{
				Params: mcp.GetPromptParams{
					Name: "timezone_conversion_guide",
					Arguments: map[string]string{
						"from_timezone": tt.fromTimezone,
						"to_timezone":   tt.toTimezone,
					},
				},
			}

			result, err := TimezoneConversionGuide(ctx, request)
			if err != nil {
				t.Fatalf("TimezoneConversionGuide() error = %v", err)
			}

			if result == nil {
				t.Fatal("TimezoneConversionGuide() returned nil result")
			}

			if len(result.Messages) == 0 {
				t.Fatal("TimezoneConversionGuide() returned no messages")
			}

			message := result.Messages[0]
			if textContent, ok := mcp.AsTextContent(message.Content); ok {
				content := textContent.Text
				for _, expected := range tt.expectContains {
					if !strings.Contains(content, expected) {
						t.Errorf("TimezoneConversionGuide() content missing expected string %q", expected)
					}
				}
			} else {
				t.Fatal("TimezoneConversionGuide() did not return text content")
			}
		})
	}
}

func TestRelativeTimeExamples(t *testing.T) {
	ctx := context.Background()

	request := mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Name: "relative_time_examples",
		},
	}

	result, err := RelativeTimeExamples(ctx, request)
	if err != nil {
		t.Fatalf("RelativeTimeExamples() error = %v", err)
	}

	if result == nil {
		t.Fatal("RelativeTimeExamples() returned nil result")
	}

	if len(result.Messages) == 0 {
		t.Fatal("RelativeTimeExamples() returned no messages")
	}

	message := result.Messages[0]
	if textContent, ok := mcp.AsTextContent(message.Content); ok {
		content := textContent.Text

		expectedStrings := []string{
			"Relative Time Expression Examples",
			"yesterday",
			"tomorrow",
			"5 minutes ago",
			"next week",
			"last month",
			"now",
			"today",
		}

		for _, expected := range expectedStrings {
			if !strings.Contains(content, expected) {
				t.Errorf("RelativeTimeExamples() content missing expected string %q", expected)
			}
		}
	} else {
		t.Fatal("RelativeTimeExamples() did not return text content")
	}
}

func TestRegisterPrompts(t *testing.T) {
	// This is more of an integration test to ensure prompts can be registered without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterPrompts() panicked: %v", r)
		}
	}()

	s := NewServer("test-server", "1.0.0")
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// Server initialization already calls RegisterPrompts via RegisterHandlers
	// Just verify the server was created successfully
}
