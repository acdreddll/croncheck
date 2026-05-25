package audit

import "github.com/user/croncheck/internal/parser"

// Report holds the full audit result for a set of cron entries.
type Report struct {
	Entries          []parser.CronEntry
	Conflicts        []Conflict
	Misconfigurations []Misconfiguration
}

// HasIssues returns true if the report contains any conflicts or misconfigurations.
func (r *Report) HasIssues() bool {
	return len(r.Conflicts) > 0 || len(r.Misconfigurations) > 0
}

// Run performs a full audit on the provided cron entries and returns a Report.
func Run(entries []parser.CronEntry) Report {
	return Report{
		Entries:           entries,
		Conflicts:         DetectConflicts(entries),
		Misconfigurations: DetectMisconfigurations(entries),
	}
}
