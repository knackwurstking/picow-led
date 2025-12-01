# picow-led

A web-based control interface for managing Picow LED devices.

## Features

- Web-based UI for controlling LED devices (PicoW)
- Device and group management
- Color control for individual devices
- SQLite database for persistent storage

## Architecture

This project is built with:
- Go (Golang) for the backend
- Echo web framework
- SQLite for data persistence
- Templ for HTML templating
- HTMX for dynamic UI updates

## Getting Started

### Prerequisites

- Go 1.25+
- Make

### Building

```bash
make init
make build
```

### Running

```bash
make run
```

This will start the server on `localhost:50836` with a test database.

## Project Structure

- `cmd/picow-led/` - Main application entry point
- `handlers/` - HTTP route handlers
- `services/` - Business logic and data access layers
- `models/` - Data models and interfaces
- `components/` - UI components
- `assets/` - Static assets (CSS, JS, images)

## Commands

- `make run` - Run the application
- `make build` - Build the binary
- `make test` - Run tests
- `make clean` - Clean build files

## Usage

1. Start the server: `make run`
2. Open your browser to `http://localhost:50836/picow-led`
3. Add devices and groups through the web interface
4. Control LED colors via the UI
