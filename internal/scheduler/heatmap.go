package scheduler

import (
	"fmt"
	"strings"

	"github.com/croncheck/internal/audit"
)

// HeatmapCell represents the activity level for a single minute-of-day slot.
type HeatmapCell struct {
	Minute int // 0–1439
	Hour   int
	Min    int
	Count  int // number of entries active at this minute
}

// Heatmap holds a 24×60 activity grid across all entries.
type Heatmap struct {
	Cells   [1440]HeatmapCell
	MaxLoad int // peak concurrent jobs in any single minute
}

// BuildHeatmap constructs a minute-resolution activity heatmap from audit entries.
// Each cell records how many cron entries are scheduled to fire at that minute.
func BuildHeatmap(entries []audit.Entry) (*Heatmap, error) {
	hm := &Heatmap{}
	for i := 0; i < 1440; i++ {
		hm.Cells[i] = HeatmapCell{
			Minute: i,
			Hour:   i / 60,
			Min:    i % 60,
		}
	}

	for _, e := range entries {
		minutes, err := collectMinutes(e.Expression)
		if err != nil {
			return nil, fmt.Errorf("heatmap: entry %q: %w", e.Label, err)
		}
		for _, m := range minutes {
			hm.Cells[m].Count++
			if hm.Cells[m].Count > hm.MaxLoad {
				hm.MaxLoad = hm.Cells[m].Count
			}
		}
	}
	return hm, nil
}

// HotSpots returns cells whose load is at or above the given threshold.
func (h *Heatmap) HotSpots(threshold int) []HeatmapCell {
	var out []HeatmapCell
	for _, c := range h.Cells {
		if c.Count >= threshold {
			out = append(out, c)
		}
	}
	return out
}

// ASCIIRow renders a single hour row as a compact ASCII bar (60 chars).
func (h *Heatmap) ASCIIRow(hour int) string {
	var sb strings.Builder
	for m := 0; m < 60; m++ {
		c := h.Cells[hour*60+m].Count
		switch {
		case c == 0:
			sb.WriteByte('.')
		case c < 3:
			sb.WriteByte('+')
		case c < 6:
			sb.WriteByte('#')
		default:
			sb.WriteByte('!')
		}
	}
	return sb.String()
}
