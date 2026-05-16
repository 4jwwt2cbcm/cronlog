package retention

import (
	"sort"
	"time"

	"github.com/cronlog/cronlog/internal/logentry"
)

// Pruner applies a retention Policy to a slice of log entries.
type Pruner struct {
	policy Policy
	now    func() time.Time
}

// NewPruner creates a Pruner with the given policy.
func NewPruner(p Policy) *Pruner {
	return &Pruner{
		policy: p,
		now:    time.Now,
	}
}

// Prune removes entries that violate the retention policy and returns the
// surviving entries. The input slice is not modified.
func (pr *Pruner) Prune(entries []logentry.Entry) []logentry.Entry {
	result := make([]logentry.Entry, 0, len(entries))

	for _, e := range entries {
		if pr.policy.KeepErrors && e.IsError() {
			result = append(result, e)
			continue
		}
		if !pr.policy.IsExpired(e.Timestamp) {
			result = append(result, e)
		}
	}

	// Sort by timestamp descending (newest first) before applying MaxEntries.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	if pr.policy.MaxEntries > 0 && len(result) > pr.policy.MaxEntries {
		result = result[:pr.policy.MaxEntries]
	}

	return result
}
