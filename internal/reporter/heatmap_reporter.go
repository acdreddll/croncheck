package reporter

import (
	"fmt"
	"io"

	"github.com/croncheck/internal/scheduler"
)

// WriteHeatmap renders a text heatmap of cron activity to w.
// Each row represents one hour (00–23) and each column one minute (00–59).
// Legend: '.' = no jobs, '+' = 1-2 jobs, '#' = 3-5 jobs, '!' = 6+ jobs.
func WriteHeatmap(w io.Writer, hm *scheduler.Heatmap) error {
	_, err := fmt.Fprintln(w, "Cron Activity Heatmap (hour × minute)")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, "Legend: '.' none  '+' 1-2  '#' 3-5  '!' 6+")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "     %s\n", minuteHeader())
	if err != nil {
		return err
	}

	for h := 0; h < 24; h++ {
		_, err = fmt.Fprintf(w, "%02d:00 %s\n", h, hm.ASCIIRow(h))
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "\nPeak concurrent jobs: %d\n", hm.MaxLoad)
	if err != nil {
		return err
	}

	threshold := 3
	if hm.MaxLoad >= threshold {
		hots := hm.HotSpots(threshold)
		_, err = fmt.Fprintf(w, "Hot spots (load >= %d): %d minute(s)\n", threshold, len(hots))
		if err != nil {
			return err
		}
		for _, c := range hots {
			_, err = fmt.Fprintf(w, "  %02d:%02d — %d jobs\n", c.Hour, c.Min, c.Count)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// minuteHeader returns a compact ruler string for the 60-column minute axis.
func minuteHeader() string {
	b := make([]byte, 60)
	for i := range b {
		if i%10 == 0 {
			b[i] = byte('0' + (i/10)%10)
		} else if i%5 == 0 {
			b[i] = '+'
		} else {
			b[i] = '-'
		}
	}
	return string(b)
}
