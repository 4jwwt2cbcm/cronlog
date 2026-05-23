// Package cli provides the command-line interface for cronlog.
//
// # Watch Command
//
// The watch command polls the log store at a configurable interval and
// prints any new entries as they arrive, providing a live tail experience
// for cron job output.
//
// Usage:
//
//	cronlog watch [flags]
//
// Flags:
//
//	-i, --interval duration   Polling interval (default 5s)
//	-j, --job string          Filter output to a specific job name
//	-e, --errors              Show only entries where the exit code is non-zero
//	-f, --format string       Output format: text, json, csv, table (default "text")
//
// Examples:
//
//	# Watch all jobs every 2 seconds
//	cronlog watch --interval 2s
//
//	# Watch only failures for the backup job
//	cronlog watch --job backup --errors
//
//	# Stream new entries as JSON
//	cronlog watch --format json
package cli
