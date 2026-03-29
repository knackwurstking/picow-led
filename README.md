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
./picow-led            # Start the server (default: :8080)
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

## API Documentation

### Device Control Endpoints

#### Set Device Color (RGB)

```
POST /api/devices/{id}/color?color=R,G,B
```

- Parameters: `color` (comma-separated RGB values 0-255)
- Returns: `204 No Content`

#### Set Device White Channel

```
POST /api/devices/{id}/white?white=VALUE
```

- Parameters: `white` (0-255)
- Returns: `204 No Content`

#### Set Device RGBW

```
POST /api/devices/{id}/rgbw?color=R,G,B&white=VALUE
```

- Parameters: `color` (RGB), `white` (0-255)
- Returns: `204 No Content`

### Web Interface

- **Home Page** (`/`) - Dashboard with devices, groups, and scenes
- **Device Page** (`/device`) - Individual device configuration

## Troubleshooting

### Database Issues

If the database becomes corrupted:

```bash
rm ~/.picow-led/picow-led.sqlite  # Delete and restart to recreate
```

### Permission Errors

Ensure the application has write permissions:

```bash
chmod -R 755 ~/.picow-led
```

### Server Not Starting

Check the port is available:

```bash
lsof -i :8080  # Check if port is in use
```

## License

MIT
