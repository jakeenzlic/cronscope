// Package exporter provides functionality to export cron job statistics
// and analysis results to various output formats.
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/yourorg/cronscope/internal/aggregator"
)

// Report represents a full exportable snapshot of cron job analysis.
type Report struct {
	GeneratedAt    time.Time                    `json:"generated_at"`
	Stats          []aggregator.JobStats        `json:"stats"`
	Streaks        []aggregator.FailureStreak   `json:"failure_streaks"`
	Trend          []aggregator.DaySummary      `json:"trend"`
}

// ExportJSON serialises the Report as indented JSON and writes it to w.
// Returns an error if marshalling or writing fails.
func ExportJSON(w io.Writer, report Report) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("exporter: marshal report: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("exporter: write report: %w", err)
	}
	return nil
}

// BuildReport assembles a Report from pre-computed aggregation results.
func BuildReport(
	stats []aggregator.JobStats,
	streaks []aggregator.FailureStreak,
	trend []aggregator.DaySummary,
) Report {
	return Report{
		GeneratedAt: time.Now().UTC(),
		Stats:       stats,
		Streaks:     streaks,
		Trend:       trend,
	}
}
