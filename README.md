# picow-led

A web-based control panel for Raspberry Pi Pico W LED projects.

## Features

- Device management with color control
- Group organization for multiple devices
- Scene presets for quick configurations

## Installation

### Prerequisites

- Go 1.25+
- Node.js (for frontend assets)

### Setup

```bash
make init          # Install dependencies and generate assets
make build         # Build the application
```

### macOS

```bash
make macos-install # Install with macOS-specific setup
make macos-update  # Update to latest version
```

### Linux/Windows

See Makefile for platform-specific commands.

## Usage

```bash
./picow-led -version  # Check version
./picow-led            # Start the server (default: :50835)
```

## Configuration

Environment variables:

| Variable             | Description            | Default        |
|----------------------|------------------------|----------------|
| `DB_PATH`            | Database location      | `~/.picow-led` |
| `SERVER_ADDRESS`     | Bind address           | `:50835`       |
| `SERVER_PATH_PREFIX` | URL path prefix        | ``             |
| `VERBOSE`            | Enable verbose logging | `true`         |

## Development

```bash
go run .     # Run development server
```

## TODO

- [x] Add API Documentation (see API.md template)
- [ ] Switch to slog using a json logger for better structured logging

- [ ] Find a way to monitor the performance by parsing logs
    ```json
    {
      "time": "2026-03-31T09:23:50.290851+02:00",
      "level": "INFO",
      "source": {
        "function": "github.com/labstack/echo/v4/middleware.RequestLogger.func1",
        "file": "/Users/knackwurstking/go/pkg/mod/github.com/labstack/echo/v4@v4.15.1/middleware/request_logger.go",
        "line": 280
      },
      "msg": "REQUEST",
      "method": "GET",
      "uri": "/css/output.css/",
      "status": 200,
      "latency": 78709,
      "host": "localhost:50888",
      "bytes_in": "",
      "bytes_out": 4359,
      "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.4 Safari/605.1.15",
      "remote_ip": "::1",
      "request_id": ""
    }
    ```

- [ ] Add Scenes management
- [ ] Improve templates structure

## License

MIT
