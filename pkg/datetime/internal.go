package datetime

import (
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/araddon/dateparse"
)

// defaultLayout is the default time format used when no other format is specified.
const defaultLayout = time.RFC3339

// GetDefaultFormat returns the default format layout string.
func GetDefaultFormat() string { return defaultLayout }

// defaultLocation is the default timezone (UTC) used when no other timezone is specified.
var defaultLocation = time.UTC

// GetDefaultTimezone returns the default timezone string.
func GetDefaultTimezone() string { return defaultLocation.String() }

// layouts provides a map of common time layout names to their format strings.
var layouts = map[string]string{
	"ANSIC":       time.ANSIC,
	"Unixdate":    time.UnixDate,
	"Rubydate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
	"DateTime":    time.DateTime,
	"DateOnly":    time.DateOnly,
	"TimeOnly":    time.TimeOnly,
}

// GetFormats returns a slice of all supported format names.
func GetFormats() []string { return slices.Collect(maps.Keys(layouts)) }

// dateTime represents a time value along with its original string representation.
type dateTime struct {
	time      time.Time
	inputTime string // The original string used to create the time.
}

// fromTime creates a new dateTime object from a time.Time object.
func fromTime(t time.Time) *dateTime {
	return &dateTime{time: t}
}

// fromString creates a new dateTime object from a string, assuming UTC if no timezone is specified.
func fromString(inputTime string) (dt *dateTime, err error) {
	return fromStringWithLocation(inputTime, nil)
}

// fromStringWithLocation creates a new dateTime object from a string and a specific location.
// If the inputTime string is empty, it defaults to the current time in UTC.
// If location is nil, it defaults to UTC.
func fromStringWithLocation(inputTime string, location *time.Location) (dt *dateTime, err error) {
	var t time.Time
	if inputTime != "" {
		if location == nil {
			location = defaultLocation
		}

		// Parse the input time string using the specified location.
		t, err = dateparse.ParseIn(inputTime, location)
		if err != nil {
			return nil, fmt.Errorf("invalid_time: Unable to parse input time: %s", inputTime)
		}
	} else {
		// Default to the current time if no input is provided.
		t = time.Now().UTC()
	}

	dt = &dateTime{
		time:      t,
		inputTime: inputTime,
	}

	return dt, nil
}

// format formats the dateTime object into a string using the specified layout and timezone.
// If format is empty, it attempts to infer the format from the original input string.
// If timezone is specified, it converts the time to that timezone.
func (dt dateTime) format(format, timezone string) (output string, err error) {
	if timezone != "" {
		// Apply the specified timezone to the time.
		location, err := time.LoadLocation(timezone)
		if err != nil {
			return "", fmt.Errorf("invalid_timezone: Invalid IANA timezone name: %s", timezone)
		}
		dt.time = dt.time.In(location)
	}

	// If a specific format is requested, use it. Otherwise, try to infer it.
	var layout string
	if format != "" {
		var ok bool
		// Check if the format is a predefined layout name.
		layout, ok = layouts[format]
		if !ok {
			// If not a predefined name, use the format string directly.
			layout = format
		}
	} else if dt.inputTime != "" {
		// If no format is provided, try to infer the format from the input time string.
		layout, err = dateparse.ParseFormat(dt.inputTime)
		if err != nil {
			return "", fmt.Errorf("invalid_format: Unable to parse format from input time: %s", dt.inputTime)
		}
	} else {
		// Fallback to the default layout if no other format can be determined.
		layout = defaultLayout
	}

	return dt.time.Format(layout), nil
}
