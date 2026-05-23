package aggregator

import (
	"sort"
	"time"

	"github.com/user/cronscope/internal/parser"
)

// DailyStats holds aggregated statistics for a single day.
type DailyStats struct {
	Date        time.Time
	TotalRuns   int
	Failures    int
	Successes   int
	FailureRate float64
}

// TrendReport holds a time-ordered series of daily statistics for a job.
type TrendReport struct {
	JobName string
	Days    []DailyStats
}

// ComputeTrend groups log entries for a given job by calendar day and returns
// a TrendReport ordered from oldest to newest day.
func ComputeTrend(jobName string, entries []parser.LogEntry) TrendReport {
	dayMap := make(map[string]*DailyStats)

	for _, e := range entries {
		if e.JobName != jobName {
			continue
		}
		day := e.StartTime.UTC().Truncate(24 * time.Hour)
		key := day.Format("2006-01-02")
		if _, ok := dayMap[key]; !ok {
			dayMap[key] = &DailyStats{Date: day}
		}
		ds := dayMap[key]
		ds.TotalRuns++
		if e.Success {
			ds.Successes++
		} else {
			ds.Failures++
		}
	}

	keys := make([]string, 0, len(dayMap))
	for k := range dayMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	days := make([]DailyStats, 0, len(keys))
	for _, k := range keys {
		ds := dayMap[k]
		if ds.TotalRuns > 0 {
			ds.FailureRate = float64(ds.Failures) / float64(ds.TotalRuns)
		}
		days = append(days, *ds)
	}

	return TrendReport{
		JobName: jobName,
		Days:    days,
	}
}
