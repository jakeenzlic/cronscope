// Package ssh provides utilities for connecting to remote servers over SSH
// and fetching cron log data for analysis by cronscope.
//
// Usage:
//
//	cfg := ssh.Config{
//		Host:    "prod-server.example.com",
//		Port:    22,
//		User:    "deploy",
//		KeyPath: "/home/user/.ssh/id_ed25519",
//	}
//
//	client, err := ssh.New(cfg)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	lines, err := ssh.FetchLogs(client, ssh.FetchOptions{
//		LogPath:     "/var/log/syslog",
//		Lines:       5000,
//		GrepPattern: "CRON",
//	})
//
// The returned lines can be passed directly to parser.ParseLog for further
// processing into aggregated statistics and failure patterns.
package ssh
