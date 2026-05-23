package aggregator

import (
	"testing"
	"time"

	"github.com/user/cronscope/internal/parser"
)

func trendEntries() []parser.LogEntry {
	base := time.Date(2024, 3, 10, 8, 0, 0, 0, time.UTC)
	return []parser.LogEntry{
		{JobName: "backup", StartTime: base, Success: true},
		{JobName: "backup", StartTime: base.Add(1 * time.Hour), Success: false},
		{JobName: "backup", StartTime: base.Add(25 * time.Hour), Success: false},
		{JobName: "backup", StartTime: base.Add(26 * time.Hour), Success: false},
		{JobName: "cleanup", StartTime: base, Success: true},
	}
}

func TestComputeTrend_DayCount(t *testing.T) {
	report := ComputeTrend("backup", trendEntries())
	if len(report.Days) != 2 {
		t.Fatalf("expected 2 days, got %d", len(report.Days))
	}
}

func TestComputeTrend_JobFilter(t *testing.T) {
	report := ComputeTrend("cleanup", trendEntries())
	if len(report.Days) != 1 {
		t.Fatalf("expected 1 day for cleanup, got %d", len(report.Days))
	}
	if report.Days[0].TotalRuns != 1 {
		t.Errorf("expected 1 run, got %d", report.Days[0].TotalRuns)
	}
}

func TestComputeTrend_FirstDayStats(t *testing.T) {
	report := ComputeTrend("backup", trendEntries())
	day0 := report.Days[0]
	if day0.TotalRuns != 2 {
		t.Errorf("expected 2 runs on day 0, got %d", day0.TotalRuns)
	}
	if day0.Failures != 1 {
		t.Errorf("expected 1 failure on day 0, got %d", day0.Failures)
	}
	if day0.FailureRate != 0.5 {
		t.Errorf("expected failure rate 0.5, got %f", day0.FailureRate)
	}
}

func TestComputeTrend_SecondDayAllFailures(t *testing.T) {
	report := ComputeTrend("backup", trendEntries())
	day1 := report.Days[1]
	if day1.FailureRate != 1.0 {
		t.Errorf("expected failure rate 1.0, got %f", day1.FailureRate)
	}
}

func TestComputeTrend_OrderedByDate(t *testing.T) {
	report := ComputeTrend("backup", trendEntries())
	for i := 1; i < len(report.Days); i++ {
		if !report.Days[i].Date.After(report.Days[i-1].Date) {
			t.Errorf("days not in ascending order at index %d", i)
		}
	}
}
