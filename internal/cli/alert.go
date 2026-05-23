package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

// AlertRule defines a threshold condition for alerting.
type AlertRule struct {
	Job            string
	MinFailures    int
	MinErrorRate   float64 // 0.0–1.0
}

// AlertResult holds the outcome of evaluating an alert rule.
type AlertResult struct {
	Rule     AlertRule
	Tripped  bool
	Message  string
}

// EvaluateAlerts checks each rule against the provided entries and returns results.
func EvaluateAlerts(entries []logentry.Entry, rules []AlertRule) []AlertResult {
	results := make([]AlertResult, 0, len(rules))
	for _, rule := range rules {
		var total, failures int
		for _, e := range entries {
			if rule.Job != "" && e.Job != rule.Job {
				continue
			}
			total++
			if e.IsError() {
				failures++
			}
		}
		var tripped bool
		var msg string
		if rule.MinFailures > 0 && failures >= rule.MinFailures {
			tripped = true
			msg = fmt.Sprintf("job %q: %d failures (threshold %d)", rule.Job, failures, rule.MinFailures)
		} else if rule.MinErrorRate > 0 && total > 0 {
			rate := float64(failures) / float64(total)
			if rate >= rule.MinErrorRate {
				tripped = true
				msg = fmt.Sprintf("job %q: error rate %.0f%% (threshold %.0f%%)",
					rule.Job, rate*100, rule.MinErrorRate*100)
			}
		}
		results = append(results, AlertResult{Rule: rule, Tripped: tripped, Message: msg})
	}
	return results
}

// PrintAlerts writes alert results to w, returning true if any were tripped.
func PrintAlerts(w io.Writer, results []AlertResult) bool {
	any := false
	for _, r := range results {
		if r.Tripped {
			any = true
			fmt.Fprintf(w, "ALERT: %s\n", r.Message)
		}
	}
	if !any {
		fmt.Fprintln(w, "No alerts tripped.")
	}
	return any
}

// buildAlertCmd constructs the 'alert' subcommand.
func buildAlertCmd(dbPath string) *cobra.Command {
	var job string
	var minFailures int
	var minErrorRate float64

	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Evaluate alert thresholds against stored log entries",
		Long:  strings.TrimSpace(alertDoc),
		RunE: func(cmd *cobra.Command, args []string) error {
			st, err := storage.New(dbPath)
			if err != nil {
				return fmt.Errorf("open storage: %w", err)
			}
			entries, err := st.All()
			if err != nil {
				return fmt.Errorf("read entries: %w", err)
			}
			rule := AlertRule{
				Job:          job,
				MinFailures:  minFailures,
				MinErrorRate: minErrorRate,
			}
			results := EvaluateAlerts(entries, []AlertRule{rule})
			tripped := PrintAlerts(os.Stdout, results)
			if tripped {
				os.Exit(1)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&job, "job", "", "filter to a specific job name")
	cmd.Flags().IntVar(&minFailures, "min-failures", 0, "alert if failure count >= this value")
	cmd.Flags().Float64Var(&minErrorRate, "min-error-rate", 0, "alert if error rate >= this value (0.0–1.0)")
	return cmd
}
