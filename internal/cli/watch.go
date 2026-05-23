package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"cronlog/internal/storage"
)

// WatchOptions holds configuration for the watch command.
type WatchOptions struct {
	Interval  time.Duration
	JobFilter string
	OnlyErrors bool
	Format    string
}

// buildWatchCmd constructs the watch subcommand which tails new log entries.
func buildWatchCmd(dbPath string) *cobra.Command {
	opts := &WatchOptions{}

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Tail new cron log entries in real time",
		Long:  "Poll the log store at a configurable interval and print new entries as they arrive.",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			fmt := opts.Format
			if fmt == "" {
				fmt = "text"
			}
			format, err := ParseFormat(fmt)
			if err != nil {
				return err
			}

			printer := NewPrinter(os.Stdout, format)
			return RunWatch(store, printer, opts)
		},
	}

	cmd.Flags().DurationVarP(&opts.Interval, "interval", "i", 5*time.Second, "Polling interval")
	cmd.Flags().StringVarP(&opts.JobFilter, "job", "j", "", "Filter by job name")
	cmd.Flags().BoolVarP(&opts.OnlyErrors, "errors", "e", false, "Show only error entries")
	cmd.Flags().StringVarP(&opts.Format, "format", "f", "text", "Output format: text, json, csv, table")

	return cmd
}

// RunWatch polls the store and prints entries newer than the last seen timestamp.
func RunWatch(store *storage.Store, printer interface{ Print(entries interface{}) error }, opts *WatchOptions) error {
	return runWatch(store, opts)
}
