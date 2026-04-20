package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestTimeFormatsResource(t *testing.T) {
	ctx := context.Background()

	request := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: "time://formats",
		},
	}

	contents, err := TimeFormatsResource(ctx, request)
	if err != nil {
		t.Fatalf("TimeFormatsResource() error = %v", err)
	}

	if len(contents) == 0 {
		t.Fatal("TimeFormatsResource() returned no contents")
	}

	textContent, ok := mcp.AsTextResourceContents(contents[0])
	if !ok {
		t.Fatal("TimeFormatsResource() did not return text resource contents")
	}

	if textContent.URI != request.Params.URI {
		t.Errorf("TimeFormatsResource() URI = %v, want %v", textContent.URI, request.Params.URI)
	}

	if textContent.MIMEType != "text/plain" {
		t.Errorf("TimeFormatsResource() MIMEType = %v, want text/plain", textContent.MIMEType)
	}

	expectedFormats := []string{
		"Available predefined time formats",
		"RFC3339",
		"Kitchen",
	}

	for _, expected := range expectedFormats {
		if !strings.Contains(textContent.Text, expected) {
			t.Errorf("TimeFormatsResource() text missing expected string %q", expected)
		}
	}
}

func TestCurrentTimeResource(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		uri        string
		wantError  bool
		expectText []string
	}{
		{
			name:       "UTC timezone",
			uri:        "time://current/UTC",
			wantError:  false,
			expectText: []string{"Current time in UTC"},
		},
		{
			name:       "America/New_York timezone",
			uri:        "time://current/America/New_York",
			wantError:  false,
			expectText: []string{"Current time in America/New_York"},
		},
		{
			name:      "invalid URI format",
			uri:       "time://current",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.ReadResourceRequest{
				Params: mcp.ReadResourceParams{
					URI: tt.uri,
				},
			}

			contents, err := CurrentTimeResource(ctx, request)

			if tt.wantError {
				if err == nil {
					t.Error("CurrentTimeResource() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("CurrentTimeResource() error = %v", err)
			}

			if len(contents) == 0 {
				t.Fatal("CurrentTimeResource() returned no contents")
			}

			textContent, ok := mcp.AsTextResourceContents(contents[0])
			if !ok {
				t.Fatal("CurrentTimeResource() did not return text resource contents")
			}

			if textContent.URI != request.Params.URI {
				t.Errorf("CurrentTimeResource() URI = %v, want %v", textContent.URI, request.Params.URI)
			}

			for _, expected := range tt.expectText {
				if !strings.Contains(textContent.Text, expected) {
					t.Errorf("CurrentTimeResource() text missing expected string %q", expected)
				}
			}
		})
	}
}

func TestTimezoneInfoResource(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		uri        string
		wantError  bool
		expectText []string
	}{
		{
			name:      "UTC timezone",
			uri:       "time://timezone-info/UTC",
			wantError: false,
			expectText: []string{
				"Timezone Information: UTC",
				"Current Time",
				"Local:",
				"UTC:",
				"IANA Name: UTC",
			},
		},
		{
			name:      "Asia/Tokyo timezone",
			uri:       "time://timezone-info/Asia/Tokyo",
			wantError: false,
			expectText: []string{
				"Timezone Information: Asia/Tokyo",
				"IANA Name: Asia/Tokyo",
			},
		},
		{
			name:      "invalid URI format",
			uri:       "time://timezone-info",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := mcp.ReadResourceRequest{
				Params: mcp.ReadResourceParams{
					URI: tt.uri,
				},
			}

			contents, err := TimezoneInfoResource(ctx, request)

			if tt.wantError {
				if err == nil {
					t.Error("TimezoneInfoResource() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("TimezoneInfoResource() error = %v", err)
			}

			if len(contents) == 0 {
				t.Fatal("TimezoneInfoResource() returned no contents")
			}

			textContent, ok := mcp.AsTextResourceContents(contents[0])
			if !ok {
				t.Fatal("TimezoneInfoResource() did not return text resource contents")
			}

			if textContent.MIMEType != "text/markdown" {
				t.Errorf("TimezoneInfoResource() MIMEType = %v, want text/markdown", textContent.MIMEType)
			}

			for _, expected := range tt.expectText {
				if !strings.Contains(textContent.Text, expected) {
					t.Errorf("TimezoneInfoResource() text missing expected string %q", expected)
				}
			}
		})
	}
}

func TestPopularTimezonesResource(t *testing.T) {
	ctx := context.Background()

	request := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: "time://timezones/popular",
		},
	}

	contents, err := PopularTimezonesResource(ctx, request)
	if err != nil {
		t.Fatalf("PopularTimezonesResource() error = %v", err)
	}

	if len(contents) == 0 {
		t.Fatal("PopularTimezonesResource() returned no contents")
	}

	textContent, ok := mcp.AsTextResourceContents(contents[0])
	if !ok {
		t.Fatal("PopularTimezonesResource() did not return text resource contents")
	}

	expectedTimezones := []string{
		"Popular Timezones",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
		"Australia/Sydney",
		"UTC",
	}

	for _, expected := range expectedTimezones {
		if !strings.Contains(textContent.Text, expected) {
			t.Errorf("PopularTimezonesResource() text missing expected timezone %q", expected)
		}
	}
}

func TestRelativeTimeGuideResource(t *testing.T) {
	ctx := context.Background()

	request := mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: "time://guide/relative-expressions",
		},
	}

	contents, err := RelativeTimeGuideResource(ctx, request)
	if err != nil {
		t.Fatalf("RelativeTimeGuideResource() error = %v", err)
	}

	if len(contents) == 0 {
		t.Fatal("RelativeTimeGuideResource() returned no contents")
	}

	textContent, ok := mcp.AsTextResourceContents(contents[0])
	if !ok {
		t.Fatal("RelativeTimeGuideResource() did not return text resource contents")
	}

	expectedSections := []string{
		"Relative Time Expressions Guide",
		"Current Time",
		"Past Expressions",
		"Future Expressions",
		"Specific Times",
		"yesterday",
		"tomorrow",
		"5 minutes ago",
		"next week",
		"last month",
	}

	for _, expected := range expectedSections {
		if !strings.Contains(textContent.Text, expected) {
			t.Errorf("RelativeTimeGuideResource() text missing expected section %q", expected)
		}
	}
}

func TestRegisterResources(t *testing.T) {
	// Integration test to ensure resources can be registered without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RegisterResources() panicked: %v", r)
		}
	}()

	s := NewServer("test-server", "1.0.0")
	if s == nil {
		t.Fatal("NewServer() returned nil")
	}

	// Server initialization already calls RegisterResources via RegisterHandlers
}
