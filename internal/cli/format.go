package cli

import (
	"fmt"
	"strings"
)

// OutputFormat represents a supported output format.
type OutputFormat int

const (
	FormatText OutputFormat = iota
	FormatJSON
	FormatCSV
	FormatTable
)

// ValidFormats lists all accepted format name strings.
var ValidFormats = []string{"text", "json", "csv", "table"}

// ParseFormat converts a string to an OutputFormat.
// Returns an error if the format is not recognised.
func ParseFormat(s string) (OutputFormat, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	case "table":
		return FormatTable, nil
	default:
		return FormatText, fmt.Errorf("unknown format %q: must be one of %s",
			s, strings.Join(ValidFormats, ", "))
	}
}

// String returns the canonical name of the OutputFormat.
func (f OutputFormat) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatCSV:
		return "csv"
	case FormatTable:
		return "table"
	default:
		return "text"
	}
}
