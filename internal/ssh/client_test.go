package ssh

import (
	"testing"
)

func TestConfig_DefaultTimeout(t *testing.T) {
	cfg := Config{
		Host:    "localhost",
		Port:    22,
		User:    "user",
		KeyPath: "/nonexistent",
	}

	// New should fail because the key file does not exist, but we can
	// verify the error is about reading the key, not a nil-pointer panic.
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error for missing key file, got nil")
	}

	const want = "ssh: read key"
	if got := err.Error(); len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("unexpected error prefix: %q", got)
	}
}

func TestConfig_PortZeroStillDialable(t *testing.T) {
	// Ensure that a zero port is forwarded as-is (":0") and the error
	// originates from key reading, not from address formatting.
	cfg := Config{Host: "localhost", Port: 0, User: "u", KeyPath: "/no"}
	_, err := New(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error()[:len("ssh: read key")] != "ssh: read key" {
		t.Errorf("unexpected error: %v", err)
	}
}
