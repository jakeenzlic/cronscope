package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// ExportCSV writes the report as CSV rows to w.
// Each row represents one job's aggregated statistics.
func ExportCSV(w io.Writer, report Report) error {
	cw := csv.NewWriter(w)

	header := []string{
		"job_name",
		"total_runs",
		"failures",
		"failure_rate_pct",
		"avg_duration_ms",
		"last_run",
		"last_status",
	}
	if err := cw.Write(header); err != nil {
		return fmt.Errorf("csv: write header: %w", err)
	}

	for _, s := range report.Stats {
		lastRun := ""
		if !s.LastRun.IsZero() {
			lastRun = s.LastRun.UTC().Format("2006-01-02T15:04:05Z")
		}
		row := []string{
			s.JobName,
			strconv.Itoa(s.TotalRuns),
			strconv.Itoa(s.Failures),
			strconv.FormatFloat(s.FailureRate*100, 'f', 2, 64),
			strconv.FormatFloat(float64(s.AvgDuration.Milliseconds()), 'f', 0, 64),
			lastRun,
			s.LastStatus,
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("csv: write row for %q: %w", s.JobName, err)
		}
	}

	cw.Flush()
	return cw.Error()
}
