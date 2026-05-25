package audit_test

import (
	"testing"

	"github.com/user/croncheck/internal/audit"
	"github.com/user/croncheck/internal/parser"
)

func entry(raw, min, hour, dom, month, dow string) parser.CronEntry {
	return parser.CronEntry{Raw: raw, Minute: min, Hour: hour, DayOfMonth: dom, Month: month, DayOfWeek: dow}
}

func TestDetectConflicts_Duplicate(t *testing.T) {
	entries := []parser.CronEntry{
		entry("* * * * *", "*", "*", "*", "*", "*"),
		entry("* * * * *", "*", "*", "*", "*", "*"),
	}
	conflicts := audit.DetectConflicts(entries)
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].Type != audit.ConflictDuplicate {
		t.Errorf("expected duplicate conflict, got %s", conflicts[0].Type)
	}
}

func TestDetectConflicts_NoConflict(t *testing.T) {
	entries := []parser.CronEntry{
		entry("0 * * * *", "0", "*", "*", "*", "*"),
		entry("30 * * * *", "30", "*", "*", "*", "*"),
	}
	conflicts := audit.DetectConflicts(entries)
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestDetectMisconfigurations_OutOfRange(t *testing.T) {
	entries := []parser.CronEntry{
		entry("60 * * * *", "60", "*", "*", "*", "*"),
	}
	issues := audit.DetectMisconfigurations(entries)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Field != "minute" {
		t.Errorf("expected field 'minute', got %q", issues[0].Field)
	}
}

func TestDetectMisconfigurations_ValidRange(t *testing.T) {
	entries := []parser.CronEntry{
		entry("0-30 1-23 * * *", "0-30", "1-23", "*", "*", "*"),
	}
	issues := audit.DetectMisconfigurations(entries)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestRun_ReportHasIssues(t *testing.T) {
	entries := []parser.CronEntry{
		entry("* * * * *", "*", "*", "*", "*", "*"),
		entry("* * * * *", "*", "*", "*", "*", "*"),
	}
	report := audit.Run(entries)
	if !report.HasIssues() {
		t.Error("expected report to have issues")
	}
}

func TestRun_CleanReport(t *testing.T) {
	entries := []parser.CronEntry{
		entry("0 6 * * 1", "0", "6", "*", "*", "1"),
		entry("0 18 * * 5", "0", "18", "*", "*", "5"),
	}
	report := audit.Run(entries)
	if report.HasIssues() {
		t.Errorf("expected clean report, got conflicts=%d misconfigs=%d",
			len(report.Conflicts), len(report.Misconfigurations))
	}
}
