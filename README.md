# cronscope

A terminal dashboard for monitoring and visualizing cron job execution history and failure patterns on remote servers.

---

## Installation

```bash
go install github.com/yourusername/cronscope@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronscope.git && cd cronscope && go build -o cronscope .
```

---

## Usage

Connect to a remote server and launch the dashboard:

```bash
cronscope --host user@example.com --log /var/log/syslog
```

**Common flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--host` | SSH target (`user@host`) | `localhost` |
| `--log` | Path to cron log file | `/var/log/cron` |
| `--interval` | Refresh interval in seconds | `10` |
| `--since` | Show history from N days ago | `7` |

Once running, use arrow keys to navigate jobs, `f` to filter by failure status, and `q` to quit.

---

## Features

- Real-time cron job status via SSH
- Visual failure pattern heatmap by hour and weekday
- Filter and search across job names and exit codes
- Lightweight — single binary, no agents required

---

## Requirements

- Go 1.21+
- SSH access to target servers
- Cron logs in standard syslog format

---

## License

MIT © 2024 cronscope contributors