package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/cronlog/internal/logentry"
)

// Printer writes log entries to an io.Writer in the requested format.
type Printer struct {
	w      io.Writer
	format string
}

// NewPrinter creates a Printer that writes to w using the given format ("text" or "json").
func NewPrinter(w io.Writer, format string) *Printer {
	return &Printer{w: w, format: format}
}

// Print outputs the provided entries according to the configured format.
func (p *Printer) Print(entries []logentry.Entry) error {
	switch p.format {
	case "json":
		return p.printJSON(entries)
	default:
		return p.printText(entries)
	}
}

func (p *Printer) printText(entries []logentry.Entry) error {
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "ID\tJOB\tSTATUS\tSTARTED\tDURATION"); err != nil {
		return err
	}
	for _, e := range entries {
		status := "ok"
		if e.IsError() {
			status = "error"
		}
		dur := e.FinishedAt.Sub(e.StartedAt).Round(time.Millisecond)
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n",
			e.ID,
			e.JobName,
			status,
			e.StartedAt.Format(time.RFC3339),
			dur,
		); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func (p *Printer) printJSON(entries []logentry.Entry) error {
	enc := json.NewEncoder(p.w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}
