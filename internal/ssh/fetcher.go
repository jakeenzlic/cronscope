package ssh

import (
	"fmt"
	"strings"
)

const (
	// DefaultLogPath is the standard syslog path used by most cron daemons.
	DefaultLogPath = "/var/log/syslog"
	// DefaultLines is the number of tail lines fetched when no limit is given.
	DefaultLines = 2000
)

// FetchOptions controls how log data is retrieved from the remote host.
type FetchOptions struct {
	LogPath string
	Lines   int
	// GrepPattern, when non-empty, pre-filters lines on the remote side.
	GrepPattern string
}

// FetchLogs retrieves cron log lines from a remote server via SSH.
// It uses `tail` (and optionally `grep`) to limit data transferred.
func FetchLogs(c *Client, opts FetchOptions) ([]string, error) {
	if opts.LogPath == "" {
		opts.LogPath = DefaultLogPath
	}
	if opts.Lines <= 0 {
		opts.Lines = DefaultLines
	}

	cmd := fmt.Sprintf("tail -n %d %s", opts.Lines, opts.LogPath)
	if opts.GrepPattern != "" {
		cmd = fmt.Sprintf("%s | grep -E %q", cmd, opts.GrepPattern)
	}

	raw, err := c.Run(cmd)
	if err != nil {
		return nil, fmt.Errorf("fetcher: %w", err)
	}

	lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")
	// Remove empty trailing line produced by Split.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}
