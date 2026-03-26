# Project Handover: Pico W LED Controller

## Overview
This is a **Pico W LED Controller web application** built with Go, SQLite, and modern web technologies. It provides a web-based interface for managing and controlling Pico W-based LED devices (LED strips, matrices, single LEDs).

## Architecture & Technology Stack

### Core Framework & Technologies
- **Language**: Go (v1.25.3)
- **Web Framework**: Echo v4 - Lightweight, fast HTTP server
- **Database**: SQLite 3 - Embedded database for device and group management
- **Templating Engine**: Templ UI - Go native HTML templating with components
- **Frontend CSS**: Tailwind 4.2.2 - Utility-first CSS framework
- **Build System**: Makefile with npm/tailwind integration
- **Frontend JavaScript**: HTMX for dynamic interactions, custom JS components

### Directory Structure
```
/
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── Makefile               # Build and development commands
├── package.json           # Frontend dependencies
└── internal/              # Main application code
    ├── assets/            # Static web assets (CSS, JS, images)
    ├── env/               # Environment and configuration
    ├── handlers/          # HTTP route handlers
    ├── routes/            # Route registration and endpoints
    ├── services/          # Business logic and data services
    └── templates/         # UI components and pages (Templ files)
└── pkg/                  # Shared models and libraries
    ├── models/            # Data models
    └── picow/             # Pico W specific protocols
```

## Key Components

### 1. Main Application Flow (main.go)
- Initializes SQLite database connection
- Sets up database configuration (max connections, idle connections)
- Creates Registry with services
- Configures Echo web server with middleware:
  - Request logging (with different log levels for status codes)
  - CORS support
  - Cache middleware
- Registers routes and starts server on defined port

### 2. Environment Configuration (internal/env/)
- **Server Configuration**: Starts on `:50835` by default
- **Verbose Logging**: Enabled by default (VERBOSE=true)
- **Database Path**: Configurable via DB_PATH environment variable
- **Server Path Prefix**: Configurable via SERVER_PATH_PREFIX

### 3. Services Layer (internal/services/)
Three main services compose the application logic:

#### Device Service
- Manages Pico W LED devices (LED strips, matrices, single LEDs)
- Database table: `devices`
  - id: INTEGER (PRIMARY KEY AUTOINCREMENT)
  - addr: TEXT UNIQUE NOT NULL (device address/URL)
  - name: TEXT NOT NULL
  - type: TEXT NOT NULL (led_strip, matrix, single_led)
  - duty: TEXT (JSON serialized device state/duty cycle)
- Key operations:
  - CreateTable() - Creates database table
  - Get(deviceID), GetByAddr(addr) - Retrieve devices
  - List() - Get all devices
  - Add(device), Update(device) - CRUD operations
  - Delete(deviceID) - Remove device
- Features:
  - Device validation (name, address format, type)
  - Address must contain protocol (http:// or https://)
  - Name can only contain letters, numbers, dashes, underscores
  - JSON duty cycle management with flexible parsing

#### Group Service
- Manages groups of devices for batch operations
- Database table: `groups`
  - id: INTEGER (PRIMARY KEY AUTOINCREMENT)
  - name: TEXT NOT NULL
- Similar CRUD operations as DeviceService

#### Registry Service
- Central service container that holds all services
- Initializes database tables for all services
- Provides unified access to DeviceService and GroupService

### 4. Routing (internal/routes/)
Three main route groups:

**API Endpoints** (`/api`):
- Device control commands (color, white, brightness, rgbw)
- POST endpoints for device state management

**UI Pages** (`/`):
- Home page - Dashboard with devices/groups overview
- Device detail page

**HTMX Endpoints** (`/htmx`):
- Dynamic UI components
- Device list management
- Group power controls
- Dialogs for adding/editing devices and groups
- CRUD operations via HTMX POST/DELETE methods

### 5. Frontend & UI
- **Templ Components**: Modular UI components using Templ library
- **Tailwind CSS**: Styling with utility classes compiled from input.css
- **HTMX Integration**: Dynamic content loading without full page reloads
- **JavaScript Libraries**: Dialog handling, form validation, etc.

## Build & Development Process

### Dependencies Setup
```bash
# Install Go dependencies
go mod tidy

# Install Node.js dependencies (for Tailwind CSS)
npm install --save-dev tailwindcss postcss autoprefixer
```

### Development Workflow
```bash
# Install development dependencies and generate assets
make init

# Generate assets (Templ + Tailwind CSS)
make generate

# Run development server with file watcher
gow -e=go,js,css -r run .
```

### Production Build
```bash
# Generate production assets and build binary
make build

# Output: bin/picow-led (executable)
```

## Database Schema

### Devices Table
```sql
CREATE TABLE IF NOT EXISTS devices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    addr TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    duty TEXT
);
```

### Groups Table
```sql
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);
```

## Configuration Options

### Environment Variables
- `DB_PATH`: Path to database directory (default: project root)
- `SERVER_ADDRESS`: Server listening address (default: :50835)
- `SERVER_PATH_PREFIX`: URL path prefix for routes
- `VERBOSE`: Enable verbose logging (default: true)

### Database Configuration
- Max Open Connections: 25
- Max Idle Connections: 25
- Connection Max Lifetime: Unlimited (0)

## Key Features

1. **Device Management**: Add, edit, delete Pico W LED devices
2. **Group Management**: Organize devices into groups for batch control
3. **Dynamic UI**: Real-time updates via HTMX
4. **Device Types**: Support for LED strips, matrices, and single LEDs
5. **State Management**: JSON-based duty cycle/state serialization
6. **Validation**: Comprehensive input validation for all operations
7. **Responsive Design**: Mobile-friendly interface with Tailwind CSS

## Development Notes

### Templ Components
- All UI components are `.templ` files that compile to Go code
- Component structure: `internal/templates/components/`
- Pages: `internal/templates/pages/`
- Dialogs: Modal forms for CRUD operations

### HTMX Integration
- POST methods for create/update operations
- DELETE methods for removal
- Dialog-based editing with form submissions
- Real-time device/group list updates

### Error Handling
- Custom error types in `internal/services/errors.go`
- Structured logging with different severity levels
- HTTP error handling in Echo middleware

## Dependencies & Versions

### Go Modules
```
github.com/Oudwins/tailwind-merge-go v0.2.1
github.com/a-h/templ v0.3.977
github.com/knackwurstking/ui v1.1.2-0.20260310092816-ee77254cdb8d
github.com/labstack/echo/v4 v4.15.1
github.com/mattn/go-sqlite3 v1.14.34
```

### JavaScript Dependencies
```json
"autoprefixer": "^10.4.27"
"postcss": "^8.5.8"
"tailwindcss": "^4.2.2"
```

## Running the Application

### Development Mode
```bash
# Install dependencies and generate assets
make init

# Run with development server and file watching
gow -e=go,js,css -r run .
```

### Production Mode
```bash
# Build production binary
make build

# Run the application
./bin/picow-led
```

### Database Location
- Default: `./data/picow-led.sqlite` (relative to DB_PATH)
- Database is automatically created on first run
- Tables are created by the Registry initialization

## Deployment Considerations

1. **Database**: SQLite file should be writable by the application user
2. **Port**: Default is 50835, can be changed via SERVER_ADDRESS env var
3. **Verbose**: Disable verbose logging in production with VERBOSE=false
4. **Path Prefix**: Can be configured for reverse proxy setups with SERVER_PATH_PREFIX
5. **Static Files**: Built assets are served from internal/assets/public/

## Troubleshooting & Common Issues

### Database Connection Issues
- Ensure DB_PATH directory exists and is writable
- Check file permissions on the SQLite database file
- Verify no other process has locked the database file

### Build Issues
- Run `make init` to install all dependencies
- Ensure Node.js is installed for Tailwind CSS compilation
- Use `go mod tidy` to sync Go dependencies

### Development Server
- If file watching doesn't work, check gow installation
- Verify no other process is using port 50835
- Check that all asset generation completed successfully with `make generate`

## Future Enhancements (from TODO.md)
- [ ] Add pins configuration section above device page color selection
- [x] Home Page: Store and restore active tab (devices, groups, scenes)
- [ ] Device pages: paddings, margins, and more

### Additional Future Enhancements to Consider:
- Scene management for predefined lighting configurations
- Scheduling system for automated device control
- API documentation and OpenAPI/Swagger integration
- Rate limiting for API endpoints
- Enhanced error handling and user feedback mechanisms
- Internationalization (i18n) support for multi-language UI
- User authentication and authorization
- Audit logging for device state changes

