package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"cronlog/internal/collector"
	"cronlog/internal/retention"
	"cronlog/internal/storage"
)

// Execute builds and runs the root CLI command.
func Execute() {
	root := buildRoot()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func buildRoot() *cobra.Command {
	var (
		dbPath  string
		format  string
	)

	root := &cobra.Command{
		Use:   "cronlog",
		Short: "Structured log aggregator for cron jobs",
	}

	root.PersistentFlags().StringVar(&dbPath, "db", "cronlog.json", "path to log database file")
	root.PersistentFlags().StringVar(&format, "format", "text", "output format: text or json")

	root.AddCommand(buildRunCmd(&dbPath, &format))
	root.AddCommand(buildLogsCmd(&dbPath, &format))
	root.AddCommand(buildPruneCmd(&dbPath))

	return root
}

func buildRunCmd(dbPath, format *string) *cobra.Command {
	var jobName string

	cmd := &cobra.Command{
		Use:   "run [-- command args...]",
		Short: "Run a command and record its output",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(*dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			if jobName == "" {
				jobName = args[0]
			}
			col := collector.New(store)
			entry, err := col.Run(jobName, args[0], args[1:]...)
			if err != nil {
				return fmt.Errorf("run: %w", err)
			}
			p := NewPrinter(os.Stdout, *format)
			p.Print([]interface{}{entry})
			return nil
		},
	}

	cmd.Flags().StringVar(&jobName, "job", "", "job name (defaults to command)")
	return cmd
}

func buildLogsCmd(dbPath, format *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Query stored log entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(*dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			f, err := ParseFilterFlags(cmd.Flags())
			if err != nil {
				return err
			}
			entries := store.Query(f)
			p := NewPrinter(os.Stdout, *format)
			p.Print(entries)
			return nil
		},
	}

	RegisterFilterFlags(cmd.Flags())
	return cmd
}

func buildPruneCmd(dbPath *string) *cobra.Command {
	var maxAgeDays int
	var maxEntries int
	var keepErrors bool

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove old log entries according to retention policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(*dbPath)
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			policy := retention.DefaultPolicy()
			if cmd.Flags().Changed("max-age") {
				policy.MaxAgeDays = maxAgeDays
			}
			if cmd.Flags().Changed("max-entries") {
				policy.MaxEntries = maxEntries
			}
			if cmd.Flags().Changed("keep-errors") {
				policy.KeepErrors = keepErrors
			}
			pruner := retention.NewPruner(store)
			removed, err := pruner.Prune(policy)
			if err != nil {
				return fmt.Errorf("prune: %w", err)
			}
			fmt.Fprintf(os.Stdout, "pruned %d entries\n", removed)
			return nil
		},
	}

	cmd.Flags().IntVar(&maxAgeDays, "max-age", 30, "maximum age in days")
	cmd.Flags().IntVar(&maxEntries, "max-entries", 1000, "maximum number of entries to keep")
	cmd.Flags().BoolVar(&keepErrors, "keep-errors", true, "always retain error entries")
	return cmd
}
