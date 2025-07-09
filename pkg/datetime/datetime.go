package datetime

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/tj/go-naturaldate"
)

// CurrentTime returns the current time in the specified timezone and format.
func CurrentTime(timezone, format string) (output string, err error) {
	return fromTime(time.Now()).
		format(format, timezone)
}

// ConvertTime converts a given time string from one timezone to another.
// If inputTimezone is empty, UTC is used as the default.
func ConvertTime(inputTime, inputTimezone, outputTimezone, format string) (output string, err error) {
	// Default to UTC if no input timezone is specified.
	var inputLocation = defaultLocation
	if inputTimezone != "" {
		// Load the input timezone location from the IANA timezone database.
		inputLocation, err = time.LoadLocation(inputTimezone)
		if err != nil {
			return "", fmt.Errorf("invalid_timezone: Invalid IANA input timezone name: %s", inputTimezone)
		}
	}

	dt, err := fromStringWithLocation(inputTime, inputLocation)
	if err != nil {
		return "", err
	}

	return dt.format(format, outputTimezone)
}

// TimeAdd adds a duration to a given time string and returns the result in the specified timezone and format.
func TimeAdd(inputTime, duration, timezone, format string) (output string, err error) {
	dt, err := fromString(inputTime)
	if err != nil {
		return "", err
	}

	// Parse the duration string (e.g., "2h30m").
	d, err := time.ParseDuration(duration)
	if err != nil {
		return "", fmt.Errorf("invalid_duration: Invalid duration format: %s", duration)
	}

	dt.time = dt.time.Add(d)

	return dt.format(format, timezone)
}

// RelativeTime parses a relative time string (e.g., "2 hours ago") based on a reference time.
func RelativeTime(inputTime, relativeTime, timezone, format string) (output string, err error) {
	refTime, err := fromString(inputTime)
	if err != nil {
		return "", err
	}

	// Parse the natural language relative time string.
	t, err := naturaldate.Parse(relativeTime, refTime.time)
	if err != nil {
		return "", fmt.Errorf("invalid_relative_time: Unable to parse relative time: %s", relativeTime)
	}

	dt := fromTime(t)
	dt.inputTime = inputTime // Store the original input time for format parsing.

	return dt.format(format, timezone)
}

// CompareTime compares two time strings and returns:
//   - -1 if timeA is before timeB
//   - 0 if timeA is equal to timeB
//   - 1 if timeA is after timeB
func CompareTime(timeA, timeB string) (int, error) {
	tA, err := dateparse.ParseAny(timeA)
	if err != nil {
		return -2, fmt.Errorf("invalid_time: invalid format for timeA: %q", timeA)
	}

	tB, err := dateparse.ParseAny(timeB)
	if err != nil {
		return -2, fmt.Errorf("invalid_time: invalid format for timeB: %q", timeB)
	}

	return tA.Compare(tB), nil
}
