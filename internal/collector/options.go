package collector

import "time"

// RunOption configures the behaviour of a single Run call.
type RunOption func(*runConfig)

type runConfig struct {
	timeout time.Duration
	env     []string
}

defaultRunConfig := runConfig{
	timeout: 0, // no timeout
}

// WithTimeout sets a maximum execution duration for the job.
// If the command exceeds this duration it is killed and the entry will
// reflect a non-zero exit code.
func WithTimeout(d time.Duration) RunOption {
	return func(c *runConfig) {
		c.timeout = d
	}
}

// WithEnv appends additional environment variables (in KEY=VALUE form) to
// the process environment for the job run.
func WithEnv(vars ...string) RunOption {
	return func(c *runConfig) {
		c.env = append(c.env, vars...)
	}
}

// applyOptions merges a slice of RunOptions onto a default runConfig.
func applyOptions(opts []RunOption) runConfig {
	cfg := defaultRunConfig
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
