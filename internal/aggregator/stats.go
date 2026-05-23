package aggregator

import (
	"time"

	"github.com/user/cronscope/internal/parser"
)

// JobStats holds aggregated statistics for a single cron job.
type JobStats struct {
	JobName      string
	TotalRuns    int
	SuccessCount int
	FailureCount int
	LastRun      time.Time
	LastStatus   string
	AvgDuration  time.Duration
	FailureRate  float64
}

// Aggregate computes per-job statistics from a slice of parsed log entries.
func Aggregate(entries []parser.LogEntry) map[string]*JobStats {
	type durationAccum struct {
		total time.Duration
		count int
	}

	stats := make(map[string]*JobStats)
	durAccum := make(map[string]*durationAccum)

	for _, e := range entries {
		s, ok := stats[e.JobName]
		if !ok {
			s = &JobStats{JobName: e.JobName}
			stats[e.JobName] = s
			durAccum[e.JobName] = &durationAccum{}
		}

		s.TotalRuns++
		if e.Success {
			s.SuccessCount++
		} else {
			s.FailureCount++
		}

		if e.Timestamp.After(s.LastRun) {
			s.LastRun = e.Timestamp
			if e.Success {
				s.LastStatus = "success"
			} else {
				s.LastStatus = "failure"
			}
		}

		if e.Duration > 0 {
			da := durAccum[e.JobName]
			da.total += e.Duration
			da.count++
		}
	}

	for name, s := range stats {
		da := durAccum[name]
		if da.count > 0 {
			s.AvgDuration = da.total / time.Duration(da.count)
		}
		if s.TotalRuns > 0 {
			s.FailureRate = float64(s.FailureCount) / float64(s.TotalRuns)
		}
	}

	return stats
}
