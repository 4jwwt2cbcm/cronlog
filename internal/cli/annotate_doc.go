package cli

const annotateDoc = `Annotate adds or replaces structured key=value tags on an existing log entry.

Tags are merged with any tags already present on the entry; existing keys are
overwritten if supplied again. The entry is identified by its unique ID, which
can be found via the 'logs' or 'stats' commands.

Examples:

  # Add a single tag
  cronlog annotate abc123 --tags env=prod

  # Add multiple tags at once
  cronlog annotate abc123 --tags env=prod,team=ops,tier=critical

  # Overwrite an existing tag
  cronlog annotate abc123 --tags env=staging

Tags are persisted to the store and are immediately available for filtering
via --tag flags on the 'logs', 'stats', and 'export' commands.
`
