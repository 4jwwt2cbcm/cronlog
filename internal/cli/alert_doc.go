package cli

// alertDoc is the long-form help text for the alert subcommand.
const alertDoc = `
Evaluate alert thresholds against stored cron job log entries.

The alert command inspects log entries (optionally filtered by job name) and
checks whether configured thresholds have been exceeded. If any alert is
tripped the command prints a message and exits with status 1, making it
suitable for use in monitoring scripts or CI pipelines.

Threshold flags (at least one must be non-zero):
  --min-failures    Trip alert when the absolute failure count meets or
                    exceeds this value.
  --min-error-rate  Trip alert when the ratio of failed runs to total runs
                    meets or exceeds this fraction (e.g. 0.5 = 50%).

Examples:
  # Alert if 'backup' has 3 or more failures
  cronlog alert --job backup --min-failures 3

  # Alert if any job has an error rate of 50% or higher
  cronlog alert --min-error-rate 0.5

  # Combine with a cron wrapper
  cronlog run --job backup -- /usr/local/bin/backup.sh && cronlog alert --job backup --min-failures 1
`
