// Package collector executes cron job commands, captures their combined
// stdout/stderr output, and stores structured log entries via the storage
// package.
//
// Basic usage:
//
//	s, _ := storage.New("/var/log/cronlog/jobs.json", retention.DefaultPolicy())
//	c := collector.New(s)
//	res := c.Run("backup", "/usr/local/bin/backup.sh", "--incremental")
//	if res.Entry.IsError() {
//		log.Printf("backup failed (exit %d): %s", res.Entry.ExitCode, res.Entry.Output)
//	}
package collector
