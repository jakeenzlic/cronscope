package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/cronscope/internal/aggregator"
	"github.com/yourorg/cronscope/internal/exporter"
)

func sampleReport() exporter.Report {
	return exporter.BuildReport(
		[]aggregator.JobStats{
			{JobName: "/etc/cron.daily/backup", TotalRuns: 10, Failures: 2, AvgDurationSec: 4.5},
		},
		[]aggregator.FailureStreak{
			{JobName: "/etc/cron.daily/backup", Length: 2, Start: time.Date(2024, 1, 5, 2, 0, 0, 0, time.UTC)},
		},
		[]aggregator.DaySummary{
			{Day: "2024-01-05", TotalRuns: 5, Failures: 2},
		},
	)
}

func TestExportJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := exporter.ExportJSON(&buf, sampleReport()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !json.Valid(buf.Bytes()) {
		t.Fatalf("output is not valid JSON: %s", buf.String())
	}
}

func TestExportJSON_ContainsJobName(t *testing.T) {
	var buf bytes.Buffer
	_ = exporter.ExportJSON(&buf, sampleReport())
	if !strings.Contains(buf.String(), "/etc/cron.daily/backup") {
		t.Error("expected job name in JSON output")
	}
}

func TestExportJSON_ContainsGeneratedAt(t *testing.T) {
	var buf bytes.Buffer
	_ = exporter.ExportJSON(&buf, sampleReport())
	if !strings.Contains(buf.String(), "generated_at") {
		t.Error("expected generated_at field in JSON output")
	}
}

func TestExportJSON_WriteError(t *testing.T) {
	err := exporter.ExportJSON(&errWriter{}, sampleReport())
	if err == nil {
		t.Fatal("expected error from failing writer, got nil")
	}
}

func TestBuildReport_FieldsPopulated(t *testing.T) {
	r := sampleReport()
	if len(r.Stats) != 1 {
		t.Errorf("expected 1 stat, got %d", len(r.Stats))
	}
	if len(r.Streaks) != 1 {
		t.Errorf("expected 1 streak, got %d", len(r.Streaks))
	}
	if len(r.Trend) != 1 {
		t.Errorf("expected 1 trend day, got %d", len(r.Trend))
	}
	if r.GeneratedAt.IsZero() {
		t.Error("expected non-zero GeneratedAt")
	}
}

// errWriter always returns an error on Write.
type errWriter struct{}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, bytes.ErrTooLarge
}
