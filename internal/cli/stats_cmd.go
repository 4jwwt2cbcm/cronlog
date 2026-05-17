package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/cronlog/internal/storage"
)

// buildStatsCmd constructs the "stats" subcommand which prints
// aggregated per-job statistics from the log store.
func buildStatsCmd(dbPath string) *cobra.Command {
	var job string

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show aggregated statistics for cron jobs",
		Long: `Display per-job statistics including total runs, error count,
average runtime, and the timestamp of the most recent execution.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}

			entries, err := store.All()
			if err != nil {
				return fmt.Errorf("read entries: %w", err)
			}

			if job != "" {
				filtered := entries[:0]
				for _, e := range entries {
					if e.Job == job {
						filtered = append(filtered, e)
					}
				}
				entries = filtered
			}

			stats := ComputeStats(entries)
			PrintStats(os.Stdout, stats)
			return nil
		},
	}

	cmd.Flags().StringVar(&job, "job", "", "Filter statistics to a specific job name")
	return cmd
}
