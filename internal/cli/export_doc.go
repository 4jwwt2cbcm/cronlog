// Package cli provides the command-line interface for cronlog.
//
// # Export
//
// The export subcommand allows persisting all stored log entries to a file
// in one of three formats:
//
//   - json  — a JSON array of entry objects (default)
//   - csv   — comma-separated values with a header row
//   - text  — human-readable plain-text lines
//
// Usage:
//
//	cronlog export --out /path/to/file.json
//	cronlog export --out /path/to/file.csv  --format csv
//	cronlog export --out /path/to/file.txt  --format text
//
// Flags:
//
//	-o, --out string     destination file path (required)
//	-f, --format string  export format: json, csv, text (default "json")
//
// The export command reads all entries from the configured store and writes
// them to the destination file.  It does not modify or remove any entries;
// use the prune subcommand to apply retention policies.
package cli
