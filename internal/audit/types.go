package audit

// Conflict represents a scheduling conflict between two cron entries.
type Conflict struct {
	// Message is a human-readable description of the conflict.
	Message string
	// EntryA is the raw cron expression of the first conflicting entry.
	EntryA string
	// EntryB is the raw cron expression of the second conflicting entry.
	EntryB string
}

// Misconfiguration represents an invalid or suspicious cron entry.
type Misconfiguration struct {
	// Message is a human-readable description of the misconfiguration.
	Message string
	// Entry is the raw cron expression of the offending entry.
	Entry string
}

// Result holds the combined output of an audit run.
type Result struct {
	Conflicts        []Conflict
	Misconfigurations []Misconfiguration
}

// HasIssues returns true if the result contains any conflicts or misconfigurations.
func (r Result) HasIssues() bool {
	return len(r.Conflicts) > 0 || len(r.Misconfigurations) > 0
}
