package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/croncheck/internal/audit"
	"github.com/user/croncheck/internal/reporter"
)

func TestWriteText_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(reporter.FormatText, &buf)
	err := w.Write(reporter.Report{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No issues found.") {
		t.Errorf("expected 'No issues found.' got %q", buf.String())
	}
}

func TestWriteText_WithIssues(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(reporter.FormatText, &buf)
	err := w.Write(reporter.Report{
		Conflicts: []audit.Conflict{
			{Message: "duplicate schedule", EntryA: "0 * * * *", EntryB: "0 * * * *"},
		},
		Misconfigurations: []audit.Misconfiguration{
			{Message: "hour out of range", Entry: "0 25 * * *"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[CONFLICT]") {
		t.Errorf("expected [CONFLICT] in output, got %q", out)
	}
	if !strings.Contains(out, "[MISCONFIG]") {
		t.Errorf("expected [MISCONFIG] in output, got %q", out)
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(reporter.FormatJSON, &buf)
	err := w.Write(reporter.Report{
		Conflicts: []audit.Conflict{
			{Message: "overlap detected", EntryA: "*/5 * * * *", EntryB: "*/5 * * * *"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := result["conflicts"]; !ok {
		t.Error("expected 'conflicts' key in JSON output")
	}
	if _, ok := result["summary"]; !ok {
		t.Error("expected 'summary' key in JSON output")
	}
}

func TestWriteJSON_Summary(t *testing.T) {
	var buf bytes.Buffer
	w := reporter.NewWriter(reporter.FormatJSON, &buf)
	_ = w.Write(reporter.Report{
		Conflicts:        []audit.Conflict{{Message: "c1"}, {Message: "c2"}},
		Misconfigurations: []audit.Misconfiguration{{Message: "m1"}},
	})

	var result struct {
		Summary struct {
			TotalConflicts        int `json:"total_conflicts"`
			TotalMisconfigurations int `json:"total_misconfigurations"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Summary.TotalConflicts != 2 {
		t.Errorf("expected 2 conflicts, got %d", result.Summary.TotalConflicts)
	}
	if result.Summary.TotalMisconfigurations != 1 {
		t.Errorf("expected 1 misconfiguration, got %d", result.Summary.TotalMisconfigurations)
	}
}
