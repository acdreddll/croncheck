package parser

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	cases := []struct {
		raw     string
		minute  string
		hour    string
	}{
		{"* * * * *", "*", "*"},
		{"0 12 * * *", "0", "12"},
		{"*/5 * * * *", "*/5", "*"},
		{"0 9-17 * * 1-5", "0", "9-17"},
		{"30 6 1,15 * *", "30", "6"},
	}
	for _, tc := range cases {
		t.Run(tc.raw, func(t *testing.T) {
			expr, err := Parse(tc.raw)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if expr.Minute != tc.minute {
				t.Errorf("Minute: got %q, want %q", expr.Minute, tc.minute)
			}
			if expr.Hour != tc.hour {
				t.Errorf("Hour: got %q, want %q", expr.Hour, tc.hour)
			}
		})
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []struct {
		raw    string
		reason string
	}{
		{"* * * *", "only 4 fields"},
		{"60 * * * *", "minute out of range"},
		{"* 25 * * *", "hour out of range"},
		{"* * 0 * *", "day out of range (0)"},
		{"* * * 13 *", "month out of range"},
		{"* * * * 8", "weekday out of range"},
		{"*/0 * * * *", "step of zero"},
		{"5-2 * * * *", "inverted range"},
	}
	for _, tc := range cases {
		t.Run(tc.reason, func(t *testing.T) {
			_, err := Parse(tc.raw)
			if err == nil {
				t.Errorf("expected error for %q (%s), got nil", tc.raw, tc.reason)
			}
		})
	}
}

func TestParse_RawPreserved(t *testing.T) {
	raw := "0 0 * * 0"
	expr, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if expr.Raw != raw {
		t.Errorf("Raw: got %q, want %q", expr.Raw, raw)
	}
}
