// Package filter provides utilities for filtering cron log entries
// by job name, time range, and exit status.
package filter

import (
	"time"

	"github.com/example/cronscope/internal/parser"
)

// Options holds the criteria used to filter log entries.
type Options struct {
	// JobName filters entries to only those matching this job name.
	// An empty string disables this filter.
	JobName string

	// Since discards entries with a timestamp before this time.
	// A zero value disables this filter.
	Since time.Time

	// Until discards entries with a timestamp after this time.
	// A zero value disables this filter.
	Until time.Time

	// OnlyFailures, when true, retains only entries where Success is false.
	OnlyFailures bool
}

// Apply returns the subset of entries that satisfy all criteria in opts.
func Apply(entries []parser.LogEntry, opts Options) []parser.LogEntry {
	result := make([]parser.LogEntry, 0, len(entries))
	for _, e := range entries {
		if opts.JobName != "" && e.JobName != opts.JobName {
			continue
		}
		if !opts.Since.IsZero() && e.Timestamp.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.Timestamp.After(opts.Until) {
			continue
		}
		if opts.OnlyFailures && e.Success {
			continue
		}
		result = append(result, e)
	}
	return result
}

// JobNames returns a deduplicated, sorted list of job names present in entries.
func JobNames(entries []parser.LogEntry) []string {
	seen := make(map[string]struct{})
	for _, e := range entries {
		seen[e.JobName] = struct{}{}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sortStrings(names)
	return names
}

// sortStrings is a simple insertion sort to avoid importing sort for a small slice.
func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		key := s[i]
		j := i - 1
		for j >= 0 && s[j] > key {
			s[j+1] = s[j]
			j--
		}
		s[j+1] = key
	}
}
