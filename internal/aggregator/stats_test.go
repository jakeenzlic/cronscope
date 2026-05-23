package aggregator

import (
	"testing"
	"time"

	"github.com/user/cronscope/internal/parser"
)

func sampleEntries() []parser.LogEntry {
	now := time.Now()
	return []parser.LogEntry{
		{JobName: "backup", Success: true, Timestamp: now.Add(-2 * time.Hour), Duration: 10 * time.Second},
		{JobName: "backup", Success: false, Timestamp: now.Add(-1 * time.Hour), Duration: 5 * time.Second},
		{JobName: "backup", Success: true, Timestamp: now, Duration: 15 * time.Second},
		{JobName: "cleanup", Success: true, Timestamp: now.Add(-30 * time.Minute), Duration: 3 * time.Second},
	}
}

func TestAggregate_JobCount(t *testing.T) {
	stats := Aggregate(sampleEntries())
	if len(stats) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(stats))
	}
}

func TestAggregate_TotalRuns(t *testing.T) {
	stats := Aggregate(sampleEntries())
	if stats["backup"].TotalRuns != 3 {
		t.Errorf("expected 3 total runs for backup, got %d", stats["backup"].TotalRuns)
	}
}

func TestAggregate_FailureRate(t *testing.T) {
	stats := Aggregate(sampleEntries())
	expected := 1.0 / 3.0
	if stats["backup"].FailureRate < expected-0.001 || stats["backup"].FailureRate > expected+0.001 {
		t.Errorf("expected failure rate ~%.4f, got %.4f", expected, stats["backup"].FailureRate)
	}
}

func TestAggregate_AvgDuration(t *testing.T) {
	stats := Aggregate(sampleEntries())
	expectedAvg := (10 + 5 + 15) / 3 * int(time.Second)
	if int(stats["backup"].AvgDuration) != expectedAvg {
		t.Errorf("expected avg duration %v, got %v", time.Duration(expectedAvg), stats["backup"].AvgDuration)
	}
}

func TestAggregate_LastStatus(t *testing.T) {
	stats := Aggregate(sampleEntries())
	if stats["backup"].LastStatus != "success" {
		t.Errorf("expected last status 'success', got '%s'", stats["backup"].LastStatus)
	}
}

func TestAggregate_EmptyEntries(t *testing.T) {
	stats := Aggregate([]parser.LogEntry{})
	if len(stats) != 0 {
		t.Errorf("expected empty stats map, got %d entries", len(stats))
	}
}
