package retention

import "errors"

var (
	errNegativeMaxAge     = errors.New("retention: MaxAge must be non-negative")
	errNegativeMaxEntries = errors.New("retention: MaxEntries must be non-negative")
)
