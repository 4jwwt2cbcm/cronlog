package cli

import (
	"flag"
	"time"
)

// FilterFlags holds parsed CLI filter arguments.
type FilterFlags struct {
	Job       string
	OnlyErrors bool
	Since     time.Duration
	MaxCount  int
	Output    string
}

// ParseFilterFlags parses filter-related flags from the given FlagSet.
// It does not call flag.Parse; the caller is responsible for that.
func ParseFilterFlags(fs *flag.FlagSet) *FilterFlags {
	f := &FilterFlags{}
	fs.StringVar(&f.Job, "job", "", "Filter log entries by job name")
	fs.BoolVar(&f.OnlyErrors, "errors", false, "Show only error entries")
	fs.DurationVar(&f.Since, "since", 0, "Show entries newer than this duration (e.g. 24h)")
	fs.IntVar(&f.MaxCount, "max", 0, "Maximum number of entries to return (0 = unlimited)")
	fs.StringVar(&f.Output, "output", "text", "Output format: text or json")
	return f
}
