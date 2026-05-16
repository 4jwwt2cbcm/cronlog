package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

// ExportFormat represents the supported export file formats.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
	ExportText ExportFormat = "text"
)

// ParseExportFormat parses a string into an ExportFormat.
func ParseExportFormat(s string) (ExportFormat, error) {
	switch strings.ToLower(s) {
	case "json":
		return ExportJSON, nil
	case "csv":
		return ExportCSV, nil
	case "text":
		return ExportText, nil
	default:
		return "", fmt.Errorf("unknown export format %q: must be json, csv, or text", s)
	}
}

// ExportEntries writes log entries to the given file path using the specified format.
func ExportEntries(entries []logentry.Entry, path string, format ExportFormat) error {
	f, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("create export file: %w", err)
	}
	defer f.Close()

	switch format {
	case ExportJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(entries); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
	case ExportCSV:
		p := NewCSVPrinter(f)
		return p.PrintAll(entries)
	case ExportText:
		p := NewPrinter(f)
		return p.PrintAll(entries)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
	return nil
}

// buildExportCmd constructs the export subcommand.
func buildExportCmd(dbPath string) *cobra.Command {
	var (
		outFile string
		fmtStr  string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export log entries to a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if outFile == "" {
				return fmt.Errorf("--out flag is required")
			}
			fmt, err := ParseExportFormat(fmtStr)
			if err != nil {
				return err
			}
			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			entries, err := store.All()
			if err != nil {
				return fmt.Errorf("read entries: %w", err)
			}
			if err := ExportEntries(entries, outFile, fmt); err != nil {
				return err
			}
			cmd.Printf("Exported %d entries to %s\n", len(entries), outFile)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outFile, "out", "o", "", "destination file path (required)")
	cmd.Flags().StringVarP(&fmtStr, "format", "f", "json", "export format: json, csv, text")
	return cmd
}
