package reporter

import (
	"encoding/json"
	"io"
)

type jsonReport struct {
	Conflicts        []jsonIssue `json:"conflicts"`
	Misconfigurations []jsonIssue `json:"misconfigurations"`
	Summary          jsonSummary `json:"summary"`
}

type jsonIssue struct {
	Message  string `json:"message"`
	EntryA   string `json:"entry_a,omitempty"`
	EntryB   string `json:"entry_b,omitempty"`
	Entry    string `json:"entry,omitempty"`
}

type jsonSummary struct {
	TotalConflicts        int `json:"total_conflicts"`
	TotalMisconfigurations int `json:"total_misconfigurations"`
}

func writeJSON(out io.Writer, r Report) error {
	jr := jsonReport{
		Conflicts:        make([]jsonIssue, 0, len(r.Conflicts)),
		Misconfigurations: make([]jsonIssue, 0, len(r.Misconfigurations)),
		Summary: jsonSummary{
			TotalConflicts:        len(r.Conflicts),
			TotalMisconfigurations: len(r.Misconfigurations),
		},
	}

	for _, c := range r.Conflicts {
		jr.Conflicts = append(jr.Conflicts, jsonIssue{
			Message: c.Message,
			EntryA:  c.EntryA,
			EntryB:  c.EntryB,
		})
	}

	for _, m := range r.Misconfigurations {
		jr.Misconfigurations = append(jr.Misconfigurations, jsonIssue{
			Message: m.Message,
			Entry:   m.Entry,
		})
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(jr)
}
