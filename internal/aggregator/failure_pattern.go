package aggregator

import (
	"sort"
	"time"

	"github.com/cronscope/cronscope/internal/parser"
)

// FailureStreak represents a consecutive run of failures for a single job.
type FailureStreak struct {
	JobName string
	Length  int
	Start   time.Time
	End     time.Time
}

// DetectFailureStreaks scans log entries and returns all consecutive failure
// streaks whose length is >= minLength, grouped per job name.
func DetectFailureStreaks(entries []parser.LogEntry, minLength int) []FailureStreak {
	// Group entries by job name, preserving chronological order.
	byJob := make(map[string][]parser.LogEntry)
	for _, e := range entries {
		byJob[e.JobName] = append(byJob[e.JobName], e)
	}

	var streaks []FailureStreak

	for job, jobEntries := range byJob {
		// Sort ascending by start time.
		sort.Slice(jobEntries, func(i, j int) bool {
			return jobEntries[i].StartTime.Before(jobEntries[j].StartTime)
		})

		streaks = append(streaks, extractStreaks(job, jobEntries, minLength)...)
	}

	// Sort result streaks by start time for deterministic output.
	sort.Slice(streaks, func(i, j int) bool {
		return streaks[i].Start.Before(streaks[j].Start)
	})

	return streaks
}

func extractStreaks(job string, entries []parser.LogEntry, minLength int) []FailureStreak {
	var result []FailureStreak
	var streakStart int
	inStreak := false

	for i, e := range entries {
		if !e.Success {
			if !inStreak {
				streakStart = i
				inStreak = true
			}
		} else {
			if inStreak {
				length := i - streakStart
				if length >= minLength {
					result = append(result, FailureStreak{
						JobName: job,
						Length:  length,
						Start:   entries[streakStart].StartTime,
						End:     entries[i-1].StartTime,
					})
				}
				inStreak = false
			}
		}
	}

	// Handle streak that extends to the end of the slice.
	if inStreak {
		length := len(entries) - streakStart
		if length >= minLength {
			result = append(result, FailureStreak{
				JobName: job,
				Length:  length,
				Start:   entries[streakStart].StartTime,
				End:     entries[len(entries)-1].StartTime,
			})
		}
	}

	return result
}
