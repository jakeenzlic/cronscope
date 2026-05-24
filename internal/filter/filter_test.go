package filter_test

import (
	"testing"
	"time"

	"github.com/example/cronscope/internal/filter"
	"github.com/example/cronscope/internal/parser"
)

var base = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

var testEntries = []parser.LogEntry{
	{JobName: "backup", Timestamp: base, Success: true, ExitCode: 0, Duration: 2 * time.Second},
	{JobName: "backup", Timestamp: base.Add(1 * time.Hour), Success: false, ExitCode: 1, Duration: 3 * time.Second},
	{JobName: "cleanup", Timestamp: base.Add(2 * time.Hour), Success: true, ExitCode: 0, Duration: 1 * time.Second},
	{JobName: "cleanup", Timestamp: base.Add(3 * time.Hour), Success: false, ExitCode: 2, Duration: 4 * time.Second},
	{JobName: "report", Timestamp: base.Add(4 * time.Hour), Success: true, ExitCode: 0, Duration: 5 * time.Second},
}

func TestApply_NoFilter(t *testing.T) {
	got := filter.Apply(testEntries, filter.Options{})
	if len(got) != len(testEntries) {
		t.Fatalf("expected %d entries, got %d", len(testEntries), len(got))
	}
}

func TestApply_JobName(t *testing.T) {
	got := filter.Apply(testEntries, filter.Options{JobName: "backup"})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.JobName != "backup" {
			t.Errorf("unexpected job name %q", e.JobName)
		}
	}
}

func TestApply_OnlyFailures(t *testing.T) {
	got := filter.Apply(testEntries, filter.Options{OnlyFailures: true})
	if len(got) != 2 {
		t.Fatalf("expected 2 failure entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Success {
			t.Errorf("expected failure entry but got success for job %q", e.JobName)
		}
	}
}

func TestApply_SinceFilter(t *testing.T) {
	got := filter.Apply(testEntries, filter.Options{Since: base.Add(2 * time.Hour)})
	if len(got) != 3 {
		t.Fatalf("expected 3 entries after since, got %d", len(got))
	}
}

func TestApply_UntilFilter(t *testing.T) {
	got := filter.Apply(testEntries, filter.Options{Until: base.Add(1 * time.Hour)})
	if len(got) != 2 {
		t.Fatalf("expected 2 entries before until, got %d", len(got))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	opts := filter.Options{
		JobName:      "cleanup",
		OnlyFailures: true,
	}
	got := filter.Apply(testEntries, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got[0].ExitCode != 2 {
		t.Errorf("expected exit code 2, got %d", got[0].ExitCode)
	}
}

func TestJobNames_Sorted(t *testing.T) {
	names := filter.JobNames(testEntries)
	want := []string{"backup", "cleanup", "report"}
	if len(names) != len(want) {
		t.Fatalf("expected %d names, got %d", len(want), len(names))
	}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("names[%d]: expected %q, got %q", i, want[i], n)
		}
	}
}
