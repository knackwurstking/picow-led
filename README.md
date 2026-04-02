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
| -------------------- | ---------------------- | -------------- |
| `DB_PATH`            | Database location      | `~/.picow-led` |
| `SERVER_ADDRESS`     | Bind address           | `:50835`       |
| `SERVER_PATH_PREFIX` | URL path prefix        | ``             |
| `VERBOSE`            | Enable verbose logging | `true`         |

## Development

```bash
go run .     # Run development server
```

## TODO

### v0.2.0
- [x] Upgrade from templui v1.6.0 to v1.9.4 [v0.2.0]

### v0.3.0
- [ ] Switch to slog using [https://github.com/lmittmann/tint] handler [v0.3.0]

### v0.4.0
- [ ] Add Scenes management [v0.4.0]

## License

MIT
