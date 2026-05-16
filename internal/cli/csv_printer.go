package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/example/cronlog/internal/logentry"
)

// CSVPrinter writes log entries as comma-separated values.
type CSVPrinter struct {
	w   *csv.Writer
	out io.Writer
}

// NewCSVPrinter creates a CSVPrinter that writes to out.
// It immediately writes the header row.
func NewCSVPrinter(out io.Writer) (*CSVPrinter, error) {
	cp := &CSVPrinter{
		w:   csv.NewWriter(out),
		out: out,
	}
	header := []string{"id", "job", "started_at", "duration_ms", "exit_code", "is_error", "output"}
	if err := cp.w.Write(header); err != nil {
		return nil, fmt.Errorf("csv header: %w", err)
	}
	cp.w.Flush()
	return cp, nil
}

// Print writes a single log entry as a CSV row.
func (cp *CSVPrinter) Print(e logentry.Entry) error {
	row := []string{
		e.ID,
		e.Job,
		e.StartedAt.UTC().Format("2006-01-02T15:04:05Z"),
		strconv.FormatInt(e.DurationMs, 10),
		strconv.Itoa(e.ExitCode),
		strconv.FormatBool(e.IsError()),
		e.Output,
	}
	if err := cp.w.Write(row); err != nil {
		return fmt.Errorf("csv write: %w", err)
	}
	cp.w.Flush()
	return cp.w.Error()
}

// PrintAll writes multiple entries in order.
func (cp *CSVPrinter) PrintAll(entries []logentry.Entry) error {
	for _, e := range entries {
		if err := cp.Print(e); err != nil {
			return err
		}
	}
	return nil
}
