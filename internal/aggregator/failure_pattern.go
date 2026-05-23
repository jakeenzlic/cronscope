package aggregator

import (
	"sort"
	"time"

	"github.com/user/cronscope/internal/parser"
)

// FailureWindow represents a consecutive streak of failures for a job.
type FailureWindow struct {
	JobName   string
	Start     time.Time
	End       time.Time
	Count     int
}

// DetectFailureStreaks identifies consecutive failure runs per job.
// A streak is reported when two or more consecutive failures occur.
func DetectFailureStreaks(entries []parser.LogEntry) []FailureWindow {
	// Group entries by job, sorted by timestamp.
	byJob := make(map[string][]parser.LogEntry)
	for _, e := range entries {
		byJob[e.JobName] = append(byJob[e.JobName], e)
	}

	var windows []FailureWindow

	for job, jobEntries := range byJob {
		sort.Slice(jobEntries, func(i, j int) bool {
			return jobEntries[i].Timestamp.Before(jobEntries[j].Timestamp)
		})

		streakStart := time.Time{}
		streakCount := 0
		lastTime := time.Time{}

		for _, e := range jobEntries {
			if !e.Success {
				if streakCount == 0 {
					streakStart = e.Timestamp
				}
				streakCount++
				lastTime = e.Timestamp
			} else {
				if streakCount >= 2 {
					windows = append(windows, FailureWindow{
						JobName: job,
						Start:   streakStart,
						End:     lastTime,
						Count:   streakCount,
					})
				}
				streakCount = 0
			}
		}

		if streakCount >= 2 {
			windows = append(windows, FailureWindow{
				JobName: job,
				Start:   streakStart,
				End:     lastTime,
				Count:   streakCount,
			})
		}
	}

	return windows
}
