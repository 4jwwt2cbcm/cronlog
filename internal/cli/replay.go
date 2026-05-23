package cli

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"cronlog/internal/collector"
	"cronlog/internal/storage"
)

// buildReplayCmd constructs the replay subcommand, which re-executes a previously
// recorded cron job command by its log entry ID.
func buildReplayCmd(store *storage.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replay <id>",
		Short: "Re-run the command associated with a log entry",
		Long: `Replay looks up a stored log entry by ID and re-executes its command.

The new execution is recorded as a fresh log entry. This is useful for
debugging failed jobs without manually reconstructing the command line.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]

			entry, err := store.FindByID(id)
			if err != nil {
				return fmt.Errorf("entry %q not found: %w", id, err)
			}

			if entry.Command == "" {
				return fmt.Errorf("entry %q has no recorded command", id)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Replaying: %s\n", entry.Command)

			parts := strings.Fields(entry.Command)
			if len(parts) == 0 {
				return fmt.Errorf("command is empty")
			}

			c := collector.New(store, entry.Job, exec.Command(parts[0], parts[1:]...))
			result, err := c.Run()
			if err != nil {
				return fmt.Errorf("replay execution failed: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Replay complete. New entry ID: %s (exit code: %d)\n",
				result.ID, result.ExitCode)
			return nil
		},
	}
	return cmd
}
