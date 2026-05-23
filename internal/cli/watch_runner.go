package cli

import (
	"fmt"
	"os"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

// runWatch is the core polling loop for the watch command.
func runWatch(store *storage.Store, opts *WatchOptions) error {
	printed := map[string]bool{}

	fmtParsed, err := ParseFormat(opts.Format)
	if err != nil {
		return err
	}
	printer := NewPrinter(os.Stdout, fmtParsed)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	fmt.Fprintf(os.Stderr, "Watching for new entries (interval: %s) — Ctrl+C to stop\n", opts.Interval)

	for range ticker.C {
		entries, err := store.All()
		if err != nil {
			return fmt.Errorf("read store: %w", err)
		}

		newEntries := filterNew(entries, printed, opts)
		for _, e := range newEntries {
			if err := printer.Print([]logentry.Entry{e}); err != nil {
				return err
			}
			printed[e.ID] = true
		}
	}
	return nil
}

// filterNew returns entries not yet printed that match the watch options.
func filterNew(entries []logentry.Entry, seen map[string]bool, opts *WatchOptions) []logentry.Entry {
	var out []logentry.Entry
	for _, e := range entries {
		if seen[e.ID] {
			continue
		}
		if opts.JobFilter != "" && e.Job != opts.JobFilter {
			continue
		}
		if opts.OnlyErrors && !e.IsError() {
			continue
		}
		out = append(out, e)
	}
	return out
}
