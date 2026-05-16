package cli

import (
	"testing"
)

func TestParseFormatKnown(t *testing.T) {
	cases := []struct {
		input  string
		want   OutputFormat
		wantOK bool
	}{
		{"text", FormatText, true},
		{"TEXT", FormatText, true},
		{"json", FormatJSON, true},
		{"JSON", FormatJSON, true},
		{"csv", FormatCSV, true},
		{"CSV", FormatCSV, true},
		{"", FormatText, true},
		{"xml", FormatText, false},
		{"unknown", FormatText, false},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, ok := ParseFormat(tc.input)
			if ok != tc.wantOK {
				t.Errorf("ParseFormat(%q) ok=%v, want %v", tc.input, ok, tc.wantOK)
			}
			if got != tc.want {
				t.Errorf("ParseFormat(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestOutputFormatString(t *testing.T) {
	cases := []struct {
		fmt  OutputFormat
		want string
	}{
		{FormatText, "text"},
		{FormatJSON, "json"},
		{FormatCSV, "csv"},
	}

	for _, tc := range cases {
		if got := tc.fmt.String(); got != tc.want {
			t.Errorf("OutputFormat(%d).String() = %q, want %q", int(tc.fmt), got, tc.want)
		}
	}
}

func TestValidFormats(t *testing.T) {
	formats := ValidFormats()
	if len(formats) != 3 {
		t.Fatalf("ValidFormats() returned %d entries, want 3", len(formats))
	}
	set := make(map[string]bool)
	for _, f := range formats {
		set[f] = true
	}
	for _, expected := range []string{"text", "json", "csv"} {
		if !set[expected] {
			t.Errorf("ValidFormats() missing %q", expected)
		}
	}
}
