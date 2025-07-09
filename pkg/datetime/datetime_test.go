package datetime

import (
	"testing"
	"time"
)

// TestFromTime tests the fromTime function to ensure it correctly creates a dateTime object from a time.Time.
func TestFromTime(t *testing.T) {
	d := time.Date(2025, 6, 7, 12, 34, 56, 78, time.UTC)

	dt := fromTime(d)

	if !dt.time.Equal(d) {
		t.Errorf("expected time %s, got %s", d, dt.time)
	}

	if dt.inputTime != "" {
		t.Errorf("expected empty inputTime, got %q", dt.inputTime)
	}
}

// TestFromString tests the fromString function to ensure it correctly parses a time string.
func TestFromString(t *testing.T) {
	inputTime := "2025-06-07T12:34:56Z"
	expectedTime := time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC)

	dt, err := fromString(inputTime)
	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	if !dt.time.Equal(expectedTime) {
		t.Errorf("expected time %s, got %s", expectedTime, dt.time)
	}

	if dt.inputTime != inputTime {
		t.Errorf("expected inputTime %s, got %s", inputTime, dt.inputTime)
	}
}

// TestFromStringWithLocation tests fromStringWithLocation with various timezones.
func TestFromStringWithLocation(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Errorf("invalid location: %v", err)
	}

	tests := []struct {
		name         string
		inputTime    string
		location     *time.Location
		expectedTime time.Time
	}{
		{
			"default UTC location",
			"2025-06-07T12:34:56",
			nil,
			time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC),
		},
		{
			"explicit UTC location only",
			"2025-06-07T12:34:56",
			time.UTC,
			time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC),
		},
		{
			"explicit New York location only",
			"2025-06-07T12:34:56",
			newYork,
			time.Date(2025, 6, 7, 12, 34, 56, 00, newYork),
		},
		{
			"implicit UTC location only",
			"2025-06-07T12:34:56Z",
			nil,
			time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC),
		},
		{
			"implicit New York location only",
			"2025-06-07T12:34:56-04:00",
			nil,
			time.Date(2025, 6, 7, 12, 34, 56, 00, newYork),
		},
		{
			"implicit location precedence",
			"2025-06-07T12:34:56-04:00",
			time.UTC,
			time.Date(2025, 6, 7, 12, 34, 56, 00, newYork),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dt, err := fromStringWithLocation(test.inputTime, test.location)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if !test.expectedTime.Equal(dt.time) {
				t.Errorf("expected time %s, got %s", test.expectedTime, dt.time)
			}
		})
	}
}

// TestFormat tests the format method of the dateTime struct.
func TestFormat(t *testing.T) {
	//newYork, err := time.LoadLocation("America/New_York")
	//if err != nil {
	//	t.Errorf("invalid location: %v", err)
	//}

	mustFromString := func(inputTime string) *dateTime {
		dt, err := fromString(inputTime)
		if err != nil {
			panic(err)
		}
		return dt
	}

	tests := []struct {
		name           string
		dateTime       *dateTime
		format         string
		timezone       string
		expectedOutput string
	}{
		{
			"FromTime default",
			fromTime(time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC)),
			"",
			"",
			"2025-06-07T12:34:56Z",
		},
		{
			"FromTime with timezone",
			fromTime(time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC)),
			"",
			"America/New_York",
			"2025-06-07T08:34:56-04:00",
		},
		{
			"FromTime with format",
			fromTime(time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC)),
			"Mon Jan 2 15:04:05 MST 2006",
			"",
			"Sat Jun 7 12:34:56 UTC 2025",
		},
		{
			"FromTime with timezone and format",
			fromTime(time.Date(2025, 6, 7, 12, 34, 56, 00, time.UTC)),
			"Mon Jan 2 15:04:05 MST 2006",
			"America/New_York",
			"Sat Jun 7 08:34:56 EDT 2025",
		},
		{
			"FromString default",
			mustFromString("2025-06-07T12:34:56Z"),
			"",
			"",
			"2025-06-07T12:34:56Z",
		},
		{
			"FromString with timezone",
			mustFromString("2025-06-07T12:34:56"),
			"",
			"America/New_York",
			"2025-06-07T08:34:56",
		},
		{
			"FromString with format",
			mustFromString("2025-06-07T12:34:56Z"),
			"Mon Jan 2 15:04:05 MST 2006",
			"",
			"Sat Jun 7 12:34:56 UTC 2025",
		},
		{
			"FromString with timezone and format",
			mustFromString("2025-06-07T12:34:56"),
			"Mon Jan 2 15:04:05 MST 2006",
			"America/New_York",
			"Sat Jun 7 08:34:56 EDT 2025",
		},
		{
			"FromString input format",
			mustFromString("Sat, 07 Jun 2025 12:34:56 MDT"),
			"",
			"",
			"Sat, 07 Jun 2025 12:34:56 MDT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := test.dateTime.format(test.format, test.timezone)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if output != test.expectedOutput {
				t.Errorf("expected output %q, got %q", test.expectedOutput, output)
			}
		})
	}
}

// TestConvertTime tests the ConvertTime function.
func TestConvertTime(t *testing.T) {
	tests := []struct {
		name           string
		inputTime      string
		inputTimezone  string
		outputTimezone string
		outputFormat   string
		expectedOutput string
	}{
		{
			"Defaults",
			"2025-07-08T12:34:56Z",
			"",
			"",
			"RFC3339",
			"2025-07-08T12:34:56Z",
		},
		{
			"Input timezone",
			"2025-07-08T12:34:56",
			"America/New_York",
			"",
			"RFC3339",
			"2025-07-08T12:34:56-04:00",
		},
		{
			"Output timezone",
			"2025-07-08T12:34:56Z",
			"",
			"Europe/Paris",
			"RFC3339",
			"2025-07-08T14:34:56+02:00",
		},
		{
			"Input & output timezone",
			"2025-07-08T12:34:56",
			"America/New_York",
			"Europe/Paris",
			"RFC3339",
			"2025-07-08T18:34:56+02:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := ConvertTime(test.inputTime, test.inputTimezone, test.outputTimezone, test.outputFormat)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if output != test.expectedOutput {
				t.Errorf("expected output %q, got %q", test.expectedOutput, output)
			}
		})
	}
}

// TestTimeAdd tests the TimeAdd function.
func TestTimeAdd(t *testing.T) {
	tests := []struct {
		name           string
		inputTime      string
		duration       string
		expectedOutput string
	}{
		{
			"+1h",
			"2025-07-08T12:34:56Z",
			"1h",
			"2025-07-08T13:34:56Z",
		},
		{
			"-1h",
			"2025-07-08T12:34:56Z",
			"-1h",
			"2025-07-08T11:34:56Z",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := TimeAdd(test.inputTime, test.duration, "", "")
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if output != test.expectedOutput {
				t.Errorf("expected output %q, got %q", test.expectedOutput, output)
			}
		})
	}
}

// TestRelativeTime tests the RelativeTime function.
func TestRelativeTime(t *testing.T) {
	tests := []struct {
		name           string
		inputTime      string
		relativeTime   string
		expectedOutput string
	}{
		{
			"now",
			"2025-07-08T12:34:56Z",
			"now",
			"2025-07-08T12:34:56Z",
		},
		{
			"now CEST",
			"2025-07-08T12:34:56+02:00",
			"now",
			"2025-07-08T12:34:56+02:00",
		},
		{
			"today",
			"2025-07-08T12:34:56Z",
			"today",
			"2025-07-08T00:00:00Z",
		},
		{
			"today CEST",
			"2025-07-08T12:34:56+02:00",
			"today",
			"2025-07-08T00:00:00+02:00",
		},
		{
			"yesterday",
			"2025-07-08T12:34:56Z",
			"yesterday",
			"2025-07-07T00:00:00Z",
		},
		{
			"yesterday CEST",
			"2025-07-08T12:34:56+02:00",
			"yesterday",
			"2025-07-07T00:00:00+02:00",
		},
		{
			"tomorrow",
			"2025-07-08T12:34:56Z",
			"tomorrow",
			"2025-07-09T00:00:00Z",
		},
		{
			"tomorrow CEST",
			"2025-07-08T12:34:56+02:00",
			"tomorrow",
			"2025-07-09T00:00:00+02:00",
		},
		{
			"5 minutes ago",
			"2025-07-08T12:34:56Z",
			"5 minutes ago",
			"2025-07-08T12:29:56Z",
		},
		{
			"5 minutes ago CEST",
			"2025-07-08T12:34:56+02:00",
			"5 minutes ago",
			"2025-07-08T12:29:56+02:00",
		},
		{
			"three days ago",
			"2025-07-08T12:34:56Z",
			"three days ago",
			"2025-07-05T00:00:00Z",
		},
		{
			"three days ago CEST",
			"2025-07-08T12:34:56+02:00",
			"three days ago",
			"2025-07-05T00:00:00+02:00",
		},
		{
			"last sunday at 5:30pm",
			"2025-07-08T12:34:56Z",
			"last sunday at 5:30pm",
			"2025-07-06T17:30:00Z",
		},
		{
			"last sunday at 5:30pm CEST",
			"2025-07-08T12:34:56+02:00",
			"last sunday at 5:30pm",
			"2025-07-06T17:30:00+02:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := RelativeTime(test.inputTime, test.relativeTime, "", "")
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if output != test.expectedOutput {
				t.Errorf("expected output %q, got %q", test.expectedOutput, output)
			}
		})
	}
}

// TestCompareTime tests the CompareTime function.
func TestCompareTime(t *testing.T) {
	tests := []struct {
		name           string
		timeA          string
		timeB          string
		expectedResult int
	}{
		{
			"equal",
			"2025-07-08T12:34:56Z",
			"2025-07-08T12:34:56Z",
			0,
		},
		{
			"timeA < timeB",
			"2025-07-08T12:34:56Z",
			"2025-07-08T12:34:56.1Z",
			-1,
		},
		{
			"timeA > timeB",
			"2025-07-08T12:34:56.1Z",
			"2025-07-08T12:34:56Z",
			1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CompareTime(test.timeA, test.timeB)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			if result != test.expectedResult {
				t.Errorf("expected output %d, got %d", test.expectedResult, result)
			}
		})
	}
}
