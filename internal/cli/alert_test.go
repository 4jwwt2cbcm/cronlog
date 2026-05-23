package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"cronlog/internal/logentry"
)

func makeAlertEntry(job string, exitCode int) logentry.Entry {
	return logentry.Entry{
		ID:        "test-id",
		Job:       job,
		StartedAt: time.Now(),
		ExitCode:  exitCode,
		Output:    "",
	}
}

func TestEvaluateAlertsNoFailures(t *testing.T) {
	entries := []logentry.Entry{
		makeAlertEntry("backup", 0),
		makeAlertEntry("backup", 0),
	}
	rules := []AlertRule{{Job: "backup", MinFailures: 1}}
	results := EvaluateAlerts(entries, rules)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Tripped {
		t.Error("expected alert not to be tripped")
	}
}

func TestEvaluateAlertsMinFailuresTripped(t *testing.T) {
	entries := []logentry.Entry{
		makeAlertEntry("backup", 0),
		makeAlertEntry("backup", 1),
		makeAlertEntry("backup", 2),
	}
	rules := []AlertRule{{Job: "backup", MinFailures: 2}}
	results := EvaluateAlerts(entries, rules)
	if !results[0].Tripped {
		t.Error("expected alert to be tripped")
	}
	if !strings.Contains(results[0].Message, "backup") {
		t.Errorf("message should mention job name, got: %s", results[0].Message)
	}
}

func TestEvaluateAlertsErrorRateTripped(t *testing.T) {
	entries := []logentry.Entry{
		makeAlertEntry("sync", 0),
		makeAlertEntry("sync", 1),
	}
	rules := []AlertRule{{Job: "sync", MinErrorRate: 0.5}}
	results := EvaluateAlerts(entries, rules)
	if !results[0].Tripped {
		t.Error("expected error rate alert to be tripped")
	}
}

func TestEvaluateAlertsJobFilter(t *testing.T) {
	entries := []logentry.Entry{
		makeAlertEntry("other", 1),
		makeAlertEntry("other", 1),
		makeAlertEntry("backup", 0),
	}
	// rule targets 'backup' only — no failures for that job
	rules := []AlertRule{{Job: "backup", MinFailures: 1}}
	results := EvaluateAlerts(entries, rules)
	if results[0].Tripped {
		t.Error("alert should not trip for a different job's failures")
	}
}

func TestPrintAlertsNoTrip(t *testing.T) {
	var buf bytes.Buffer
	results := []AlertResult{{Tripped: false}}
	tripped := PrintAlerts(&buf, results)
	if tripped {
		t.Error("PrintAlerts should return false when no alerts tripped")
	}
	if !strings.Contains(buf.String(), "No alerts") {
		t.Errorf("expected 'No alerts' message, got: %s", buf.String())
	}
}

func TestPrintAlertsTripped(t *testing.T) {
	var buf bytes.Buffer
	results := []AlertResult{{Tripped: true, Message: "job \"backup\": 3 failures (threshold 2)"}}
	tripped := PrintAlerts(&buf, results)
	if !tripped {
		t.Error("PrintAlerts should return true when an alert is tripped")
	}
	if !strings.Contains(buf.String(), "ALERT:") {
		t.Errorf("expected ALERT prefix, got: %s", buf.String())
	}
}
