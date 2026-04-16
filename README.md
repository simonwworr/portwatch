# portwatch

Lightweight CLI to monitor and alert on open port changes across hosts.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Scan a host and watch for port changes:

```bash
portwatch watch --host 192.168.1.1 --interval 60s
```

Scan multiple hosts from a config file:

```bash
portwatch watch --config hosts.yaml --alert email
```

Example `hosts.yaml`:

```yaml
hosts:
  - 192.168.1.1
  - 192.168.1.2
  - example.com
interval: 30s
alert: stdout
```

When a port opens or closes, portwatch will output:

```
[ALERT] 192.168.1.1 — port 8080 OPENED (2024-01-15 10:32:01)
[ALERT] 192.168.1.2 — port 22 CLOSED  (2024-01-15 10:32:04)
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--host` | Target host to monitor | — |
| `--config` | Path to hosts config file | — |
| `--interval` | Scan interval | `60s` |
| `--alert` | Alert output (`stdout`, `email`) | `stdout` |
| `--ports` | Port range to scan | `1-65535` |

## License

MIT © [yourusername](https://github.com/yourusername)