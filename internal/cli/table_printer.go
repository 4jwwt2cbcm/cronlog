package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/example/cronlog/internal/logentry"
)

// TablePrinter renders log entries as an aligned table.
type TablePrinter struct {
	w *tabwriter.Writer
}

// NewTablePrinter creates a TablePrinter writing to w.
func NewTablePrinter(w io.Writer) *TablePrinter {
	return &TablePrinter{
		w: tabwriter.NewWriter(w, 0, 0, 2, ' ', 0),
	}
}

// PrintHeader writes the column header row.
func (p *TablePrinter) PrintHeader() {
	fmt.Fprintln(p.w, "ID\tJOB\tSTATUS\tEXIT\tSTARTED\tDURATION")
	fmt.Fprintln(p.w, strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 7)+"\t"+strings.Repeat("-", 4)+"\t"+strings.Repeat("-", 19)+"\t"+strings.Repeat("-", 10))
}

// Print writes a single log entry row.
func (p *TablePrinter) Print(e logentry.Entry) {
	status := "OK"
	if e.IsError() {
		status = "ERROR"
	}
	started := e.StartedAt.Format(time.RFC3339)
	duration := e.Duration.Round(time.Millisecond).String()
	fmt.Fprintf(p.w, "%s\t%s\t%s\t%d\t%s\t%s\n",
		e.ID[:8], e.JobName, status, e.ExitCode, started, duration)
}

// PrintAll writes the header followed by all entries and flushes.
func (p *TablePrinter) PrintAll(entries []logentry.Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(p.w, "No log entries found.")
		p.w.Flush()
		return
	}
	p.PrintHeader()
	for _, e := range entries {
		p.Print(e)
	}
	p.w.Flush()
}
