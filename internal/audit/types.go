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

// Summary returns a brief string describing the number of issues found.
func (r Result) Summary() string {
	if !r.HasIssues() {
		return "no issues found"
	}
	var parts []string
	if n := len(r.Conflicts); n > 0 {
		parts = append(parts, fmt.Sprintf("%d conflict(s)", n))
	}
	if n := len(r.Misconfigurations); n > 0 {
		parts = append(parts, fmt.Sprintf("%d misconfiguration(s)", n))
	}
	return strings.Join(parts, ", ")
}
