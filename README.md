# cronlog

Structured log aggregator for cron jobs with filtering and retention policies.

## Installation

```bash
go install github.com/yourname/cronlog@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/cronlog.git && cd cronlog && go build ./...
```

## Usage

Wrap any cron command with `cronlog` to capture and store structured output:

```bash
cronlog run --job "backup" --retain 30d -- /usr/local/bin/backup.sh
```

Filter logs by job name, status, or time range:

```bash
# View logs for a specific job
cronlog logs --job "backup" --since 7d

# Show only failed runs
cronlog logs --status failed

# Tail live output
cronlog tail --job "backup"
```

Configure retention policies in `cronlog.yaml`:

```yaml
retention:
  default: 30d
  jobs:
    backup: 90d
    cleanup: 7d
filters:
  suppress_empty: true
  min_duration: 1s
```

Run the log server to expose a query API:

```bash
cronlog server --addr :8080 --storage /var/lib/cronlog
```

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--storage` | `~/.cronlog` | Path to log storage directory |
| `--retain` | `30d` | Default log retention period |
| `--format` | `json` | Output format (`json`, `text`) |

## License

MIT © yourname