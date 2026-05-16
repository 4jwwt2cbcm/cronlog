package cli

import (
	"flag"
	"testing"
	"time"
)

func TestParseFilterFlagsDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	f := ParseFilterFlags(fs)

	if err := fs.Parse([]string{}); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if f.Job != "" {
		t.Errorf("expected empty job, got %q", f.Job)
	}
	if f.OnlyErrors {
		t.Error("expected OnlyErrors=false by default")
	}
	if f.Since != 0 {
		t.Errorf("expected Since=0, got %v", f.Since)
	}
	if f.MaxCount != 0 {
		t.Errorf("expected MaxCount=0, got %d", f.MaxCount)
	}
	if f.Output != "text" {
		t.Errorf("expected output=text, got %q", f.Output)
	}
}

func TestParseFilterFlagsValues(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	f := ParseFilterFlags(fs)

	args := []string{
		"-job", "backup",
		"-errors",
		"-since", "48h",
		"-max", "50",
		"-output", "json",
	}
	if err := fs.Parse(args); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if f.Job != "backup" {
		t.Errorf("expected job=backup, got %q", f.Job)
	}
	if !f.OnlyErrors {
		t.Error("expected OnlyErrors=true")
	}
	if f.Since != 48*time.Hour {
		t.Errorf("expected Since=48h, got %v", f.Since)
	}
	if f.MaxCount != 50 {
		t.Errorf("expected MaxCount=50, got %d", f.MaxCount)
	}
	if f.Output != "json" {
		t.Errorf("expected output=json, got %q", f.Output)
	}
}
