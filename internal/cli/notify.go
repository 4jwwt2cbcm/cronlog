package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

// NotifyConfig holds configuration for alert notifications.
type NotifyConfig struct {
	Hook       string
	MinFails   int
	ErrorRate  float64
	JobFilter  string
	DryRun     bool
}

// RunNotify evaluates alerts and dispatches a webhook/command if thresholds are met.
func RunNotify(w io.Writer, store *storage.Store, cfg NotifyConfig) error {
	entries := store.All()

	var filtered []logentry.Entry
	for _, e := range entries {
		if cfg.JobFilter == "" || e.Job == cfg.JobFilter {
			filtered = append(filtered, e)
		}
	}

	alerts := EvaluateAlerts(filtered, cfg.MinFails, cfg.ErrorRate)
	if len(alerts) == 0 {
		fmt.Fprintln(w, "no alerts triggered")
		return nil
	}

	PrintAlerts(w, alerts)

	if cfg.Hook == "" {
		return nil
	}

	var lines []string
	for _, a := range alerts {
		lines = append(lines, fmt.Sprintf("%s: %s", a.Job, a.Reason))
	}
	message := strings.Join(lines, "; ")

	if cfg.DryRun {
		fmt.Fprintf(w, "[dry-run] would invoke: %s %q\n", cfg.Hook, message)
		return nil
	}

	cmd := exec.Command(cfg.Hook, message)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("notify hook failed: %w", err)
	}
	return nil
}

// buildNotifyCmd constructs the cobra command for the notify subcommand.
func buildNotifyCmd(dbPath string) *cobra.Command {
	var cfg NotifyConfig

	cmd := &cobra.Command{
		Use:   "notify",
		Short: "Evaluate alert thresholds and invoke a notification hook",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := storage.New(dbPath)
			if err != nil {
				return err
			}
			return RunNotify(cmd.OutOrStdout(), store, cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.Hook, "hook", "", "executable to invoke with alert summary as argument")
	cmd.Flags().IntVar(&cfg.MinFails, "min-failures", 1, "minimum consecutive failures to trigger alert")
	cmd.Flags().Float64Var(&cfg.ErrorRate, "error-rate", 0.5, "error rate threshold (0.0–1.0) to trigger alert")
	cmd.Flags().StringVar(&cfg.JobFilter, "job", "", "filter alerts to a specific job")
	cmd.Flags().BoolVar(&cfg.DryRun, "dry-run", false, "print hook invocation without executing")

	return cmd
}
