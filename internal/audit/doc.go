// Package audit analyses a slice of parsed cron entries and produces a Report
// that describes any scheduling conflicts or misconfigurations found.
//
// Conflict detection identifies entries that are duplicates or that will fire
// at the same time. Misconfiguration detection checks that field values fall
// within their valid numeric ranges and that step/range syntax is well-formed.
//
// Typical usage:
//
//	report := audit.Run(entries)
//	if report.HasIssues() {
//		// handle report.Conflicts and report.Misconfigurations
//	}
package audit
