package parser

import (
	"strings"
	"testing"
	"time"
)

const sampleLog = `Jan  2 15:04:05 myhost CRON[1234]: (root) CMD (/usr/bin/backup.sh)
Jan  2 15:10:00 myhost CRON[1235]: (www-data) CMD (/opt/cleanup.sh exited with code 1)
Jan  2 15:15:30 myhost CRON[1236]: (deploy) CMD (/opt/deploy.sh exited with code 0)
this line should be ignored
`

func TestParseLog_EntryCount(t *testing.T) {
	entries, err := ParseLog(strings.NewReader(sampleLog), 2024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestParseLog_SuccessEntry(t *testing.T) {
	entries, _ := ParseLog(strings.NewReader(sampleLog), 2024)
	e := entries[0]

	if e.Status != StatusSuccess {
		t.Errorf("expected success, got %s", e.Status)
	}
	if e.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", e.ExitCode)
	}
	if e.JobName != "/usr/bin/backup.sh" {
		t.Errorf("unexpected job name: %s", e.JobName)
	}
}

func TestParseLog_FailureEntry(t *testing.T) {
	entries, _ := ParseLog(strings.NewReader(sampleLog), 2024)
	e := entries[1]

	if e.Status != StatusFailure {
		t.Errorf("expected failure, got %s", e.Status)
	}
	if e.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", e.ExitCode)
	}
}

func TestParseLog_ExitCodeZeroIsSuccess(t *testing.T) {
	entries, _ := ParseLog(strings.NewReader(sampleLog), 2024)
	e := entries[2]

	if e.Status != StatusSuccess {
		t.Errorf("expected success for exit code 0, got %s", e.Status)
	}
}

func TestParseLog_Timestamp(t *testing.T) {
	entries, _ := ParseLog(strings.NewReader(sampleLog), 2024)
	e := entries[0]

	expected := time.Date(2024, time.January, 2, 15, 4, 5, 0, time.UTC)
	if !e.Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, e.Timestamp)
	}
}

func TestParseLog_EmptyInput(t *testing.T) {
	entries, err := ParseLog(strings.NewReader(""), 2024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}
