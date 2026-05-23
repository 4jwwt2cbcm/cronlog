package cli

// notify subcommand — long-form documentation.

const notifyLongDesc = `
The notify command evaluates alert thresholds across stored cron job log
entries and, when any threshold is exceeded, invokes a user-supplied
executable (hook) with a human-readable summary of triggered alerts.

Alert thresholds
  --min-failures   Trigger when a job has at least N consecutive failures.
  --error-rate     Trigger when a job's failure rate exceeds the given ratio
                   (e.g. 0.5 means more than 50 %% of runs failed).

Hook invocation
  The hook executable receives a single string argument containing a
  semicolon-separated list of "job: reason" pairs. The hook's stdout is
  forwarded to cronlog's own stdout. A non-zero exit code from the hook
  is reported as an error.

Examples

  # Print alerts to stdout without invoking any hook
  cronlog notify --min-failures 3

  # Invoke a Slack webhook script when error rate exceeds 80 %%
  cronlog notify --error-rate 0.8 --hook ./scripts/slack-notify.sh

  # Preview what would be invoked without actually running the hook
  cronlog notify --hook ./scripts/slack-notify.sh --dry-run

  # Restrict evaluation to a single job
  cronlog notify --job backup --min-failures 2 --hook ./scripts/pagerduty.sh
`
