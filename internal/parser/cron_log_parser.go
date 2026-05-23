package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
)

// JobStatus represents the execution result of a cron job.
type JobStatus string

const (
	StatusSuccess JobStatus = "success"
	StatusFailure JobStatus = "failure"
	StatusUnknown JobStatus = "unknown"
)

// CronEntry represents a single parsed cron job execution record.
type CronEntry struct {
	Timestamp time.Time
	JobName   string
	Status    JobStatus
	ExitCode  int
	Message   string
}

// logLineRegex matches syslog-style cron entries:
// e.g. "Jan  2 15:04:05 host CRON[1234]: (user) CMD (job-name)"
var logLineRegex = regexp.MustCompile(
	`^(\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2})\s+\S+\s+CRON\[\d+\]:\s+\(\S+\)\s+(\w+)\s+\((.+)\)`,
)

// exitCodeRegex captures exit codes from messages like "job-name exited with code 1"
var exitCodeRegex = regexp.MustCompile(`exited with code (\d+)`)

// ParseLog reads cron log lines from r and returns a slice of CronEntry.
func ParseLog(r io.Reader, year int) ([]CronEntry, error) {
	var entries []CronEntry
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := parseLine(line, year)
		if err != nil {
			// skip lines that don't match
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scanning log: %w", err)
	}

	return entries, nil
}

func parseLine(line string, year int) (CronEntry, error) {
	matches := logLineRegex.FindStringSubmatch(line)
	if matches == nil {
		return CronEntry{}, fmt.Errorf("parser: no match")
	}

	rawTime := fmt.Sprintf("%d %s", year, matches[1])
	ts, err := time.Parse("2006 Jan  2 15:04:05", rawTime)
	if err != nil {
		ts, err = time.Parse("2006 Jan _2 15:04:05", rawTime)
		if err != nil {
			return CronEntry{}, fmt.Errorf("parser: parse time: %w", err)
		}
	}

	cmd := matches[2] // CMD, START, etc.
	msg := matches[3]

	status := StatusUnknown
	exitCode := 0

	if cmd == "CMD" {
		status = StatusSuccess
	}

	if codeMatch := exitCodeRegex.FindStringSubmatch(msg); codeMatch != nil {
		fmt.Sscanf(codeMatch[1], "%d", &exitCode)
		if exitCode != 0 {
			status = StatusFailure
		} else {
			status = StatusSuccess
		}
	}

	return CronEntry{
		Timestamp: ts,
		JobName:   msg,
		Status:    status,
		ExitCode:  exitCode,
		Message:   line,
	}, nil
}
