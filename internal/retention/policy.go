package retention

import (
	"time"
)

// Policy defines how long log entries should be retained.
type Policy struct {
	// MaxAge is the maximum age of a log entry before it is considered expired.
	MaxAge time.Duration
	// MaxEntries is the maximum number of entries to retain. 0 means unlimited.
	MaxEntries int
	// KeepErrors ensures error entries are never pruned regardless of age.
	KeepErrors bool
}

// DefaultPolicy returns a sensible default retention policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:     7 * 24 * time.Hour,
		MaxEntries: 10000,
		KeepErrors: true,
	}
}

// IsExpired reports whether an entry recorded at t has exceeded the policy's MaxAge.
func (p Policy) IsExpired(t time.Time) bool {
	if p.MaxAge <= 0 {
		return false
	}
	return time.Since(t) > p.MaxAge
}

// Validate returns an error string if the policy contains invalid values.
func (p Policy) Validate() error {
	if p.MaxAge < 0 {
		return errNegativeMaxAge
	}
	if p.MaxEntries < 0 {
		return errNegativeMaxEntries
	}
	return nil
}
