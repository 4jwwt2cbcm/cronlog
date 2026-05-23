package cli

// ReplayDoc describes the replay command for documentation and help generation.
const ReplayDoc = `
replay — Re-run a previously recorded cron job

USAGE
  cronlog replay <id>

DESCRIPTION
  The replay command retrieves a stored log entry by its unique ID and
  re-executes the original command in a new collector run. The resulting
  output and exit code are stored as a fresh log entry, leaving the
  original entry untouched.

  This is especially useful when investigating failures: instead of
  manually reconstructing a complex command, replay fetches all details
  from the stored entry and runs it again under the same job name.

EXAMPLES
  # Replay a specific entry
  cronlog replay 4f3a1b2c

  # Combine with logs to find the ID first
  cronlog logs --job backup --only-errors
  cronlog replay <id from above>

NOTES
  - The replayed job is stored under the same job name as the original.
  - Entries without a recorded Command field cannot be replayed.
`
