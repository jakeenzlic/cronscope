package aggregator

import (
	"testing"
	"time"

	"github.com/cronscope/cronscope/internal/parser"
)

func failureEntries() []parser.LogEntry {
	now := time.Now()
	return []parser.LogEntry{
		{JobName: "backup", Success: false, StartTime: now.Add(-6 * time.Hour)},
		{JobName: "backup", Success: false, StartTime: now.Add(-5 * time.Hour)},
		{JobName: "backup", Success: false, StartTime: now.Add(-4 * time.Hour)},
		{JobName: "backup", Success: true, StartTime: now.Add(-3 * time.Hour)},
		{JobName: "backup", Success: false, StartTime: now.Add(-2 * time.Hour)},
		{JobName: "backup", Success: false, StartTime: now.Add(-1 * time.Hour)},
		{JobName: "sync", Success: false, StartTime: now.Add(-3 * time.Hour)},
		{JobName: "sync", Success: false, StartTime: now.Add(-2 * time.Hour)},
		{JobName: "sync", Success: true, StartTime: now.Add(-1 * time.Hour)},
	}
}

func TestDetectFailureStreaks_ReturnsStreaks(t *testing.T) {
	streaks := DetectFailureStreaks(failureEntries(), 2)
	if len(streaks) == 0 {
		t.Fatal("expected at least one streak, got none")
	}
}

func TestDetectFailureStreaks_BackupLongestStreak(t *testing.T) {
	streaks := DetectFailureStreaks(failureEntries(), 2)
	var backupMax int
	for _, s := range streaks {
		if s.JobName == "backup" && s.Length > backupMax {
			backupMax = s.Length
		}
	}
	if backupMax != 3 {
		t.Errorf("expected backup longest streak=3, got %d", backupMax)
	}
}

func TestDetectFailureStreaks_MinThresholdFilters(t *testing.T) {
	streaks := DetectFailureStreaks(failureEntries(), 3)
	for _, s := range streaks {
		if s.Length < 3 {
			t.Errorf("streak length %d is below min threshold 3", s.Length)
		}
	}
}

func TestDetectFailureStreaks_StreakStartTime(t *testing.T) {
	streaks := DetectFailureStreaks(failureEntries(), 2)
	for _, s := range streaks {
		if s.Start.IsZero() {
			t.Errorf("streak for job %q has zero start time", s.JobName)
		}
	}
}

func TestDetectFailureStreaks_NoFailures(t *testing.T) {
	entries := []parser.LogEntry{
		{JobName: "clean", Success: true},
		{JobName: "clean", Success: true},
	}
	streaks := DetectFailureStreaks(entries, 2)
	if len(streaks) != 0 {
		t.Errorf("expected no streaks for all-success entries, got %d", len(streaks))
	}
}
