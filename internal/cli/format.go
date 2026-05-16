package cli

import "strings"

// OutputFormat represents the supported output formats for log display.
type OutputFormat int

const (
	FormatText OutputFormat = iota
	FormatJSON
	FormatCSV
)

// ParseFormat converts a string flag value to an OutputFormat.
// It returns FormatText and false if the format is unrecognized.
func ParseFormat(s string) (OutputFormat, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, true
	case "json":
		return FormatJSON, true
	case "csv":
		return FormatCSV, true
	default:
		return FormatText, false
	}
}

// String returns the canonical string name of the OutputFormat.
func (f OutputFormat) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatCSV:
		return "csv"
	default:
		return "text"
	}
}

// ValidFormats returns a slice of all supported format name strings.
func ValidFormats() []string {
	return []string{"text", "json", "csv"}
}
