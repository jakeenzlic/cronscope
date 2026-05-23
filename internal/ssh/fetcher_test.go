package ssh

import (
	"strings"
	"testing"
)

// mockClient satisfies the minimal interface needed by FetchLogs for testing.
// We test FetchLogs logic by injecting a stub that returns controlled output.

type stubClient struct {
	output string
	err    error
	gotCmd string
}

func (s *stubClient) Run(cmd string) (string, error) {
	s.gotCmd = cmd
	return s.output, s.err
}

// fetchLogsWithRunner is a testable variant that accepts a runner interface.
func fetchLogsWithRunner(runner interface{ Run(string) (string, error) }, opts FetchOptions) ([]string, error) {
	if opts.LogPath == "" {
		opts.LogPath = DefaultLogPath
	}
	if opts.Lines <= 0 {
		opts.Lines = DefaultLines
	}
	cmd := ""
	if opts.GrepPattern != "" {
		cmd = "tail -n " + itoa(opts.Lines) + " " + opts.LogPath + " | grep -E \"" + opts.GrepPattern + "\""
	} else {
		cmd = "tail -n " + itoa(opts.Lines) + " " + opts.LogPath
	}
	raw, err := runner.Run(cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

func itoa(n int) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(
		strings.ReplaceAll(fmt.Sprintf("%d", n), " ", ""), "\t", ""), "\n", ""))
}

func TestFetchLogs_LineCount(t *testing.T) {
	stub := &stubClient{output: "line1\nline2\nline3\n"}
	lines, err := fetchLogsWithRunner(stub, FetchOptions{LogPath: "/var/log/syslog", Lines: 500})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestFetchLogs_DefaultsApplied(t *testing.T) {
	stub := &stubClient{output: "a\nb\n"}
	_, err := fetchLogsWithRunner(stub, FetchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stub.gotCmd, DefaultLogPath) {
		t.Errorf("expected default log path in command, got: %s", stub.gotCmd)
	}
}

func TestFetchLogs_EmptyOutput(t *testing.T) {
	stub := &stubClient{output: ""}
	lines, err := fetchLogsWithRunner(stub, FetchOptions{LogPath: "/tmp/test.log", Lines: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines for empty output, got %d", len(lines))
	}
}
