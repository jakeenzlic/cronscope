// Package filter provides in-memory filtering of parsed cron log entries.
//
// Entries can be narrowed by job name, time window, or success/failure status.
// Filters are composable: all non-zero criteria in an Options struct are applied
// together as a logical AND.
//
// Example usage:
//
//	entries, _ := parser.ParseLog(r)
//	filtered := filter.Apply(entries, filter.Options{
//		JobName:      "backup",
//		Since:        time.Now().Add(-24 * time.Hour),
//		OnlyFailures: true,
//	})
//
// The JobNames helper returns a deduplicated sorted list of job names, which is
// useful for populating selection widgets in the terminal dashboard.
package filter
