package exporter

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/user/cronscope/internal/aggregator"
)

func TestExportCSV_HeaderRow(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(&buf, sampleReport()); err != nil {
		t.Fatalf("ExportCSV returned error: %v", err)
	}

	r := csv.NewReader(&buf)
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv parse error: %v", err)
	}
	if len(rows) == 0 {
		t.Fatal("expected at least one row")
	}
	if rows[0][0] != "job_name" {
		t.Errorf("expected header 'job_name', got %q", rows[0][0])
	}
}

func TestExportCSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	rep := sampleReport()
	if err := ExportCSV(&buf, rep); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, _ := r.ReadAll()
	// 1 header + number of stats entries
	want := 1 + len(rep.Stats)
	if len(rows) != want {
		t.Errorf("row count: got %d, want %d", len(rows), want)
	}
}

func TestExportCSV_ContainsJobName(t *testing.T) {
	var buf bytes.Buffer
	if err := ExportCSV(&buf, sampleReport()); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	if !strings.Contains(buf.String(), "backup") {
		t.Error("expected CSV to contain job name 'backup'")
	}
}

func TestExportCSV_FailureRateFormat(t *testing.T) {
	rep := Report{
		GeneratedAt: time.Now(),
		Stats: []aggregator.JobStats{
			{JobName: "myjob", TotalRuns: 4, Failures: 1, FailureRate: 0.25},
		},
	}
	var buf bytes.Buffer
	if err := ExportCSV(&buf, rep); err != nil {
		t.Fatalf("ExportCSV error: %v", err)
	}
	if !strings.Contains(buf.String(), "25.00") {
		t.Errorf("expected failure rate '25.00' in output, got:\n%s", buf.String())
	}
}

func TestExportCSV_WriteError(t *testing.T) {
	err := ExportCSV(&errWriter{}, sampleReport())
	if err == nil {
		t.Error("expected error from failing writer, got nil")
	}
}
