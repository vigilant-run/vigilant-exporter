# Vigilant Exporter

The Vigilant Exporter is a command line tool that collects and sends data to the Vigilant API.

## Installation (Linux)

```bash
curl -L https://github.com/vigilant-io/vigilant-exporter/releases/latest/download/vigilant-exporter-linux-amd64 -o /usr/local/bin/vigilant-exporter
```

## Installation (macOS)

```bash
curl -L https://github.com/vigilant-io/vigilant-exporter/releases/latest/download/vigilant-exporter-darwin-amd64 -o /usr/local/bin/vigilant-exporter
```

## Usage

```bash
vigilant-exporter --file /path/to/log/file --token tk_1234567890 --endpoint https://ingress.vigilant.run
```

## Configuration

The Vigilant Exporter can be configured with the following flags:

- `--file`: The path to the log file to monitor.
- `--token`: The authentication token.
- `--endpoint`: The endpoint URL for log ingestion.
- `--insecure`: Send logs over HTTP instead of HTTPS.
- `--help`: Show help for any command.
