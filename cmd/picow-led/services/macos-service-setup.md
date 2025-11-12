# PicoW LED macOS Launchd Service Setup

This document explains how to set up the picow-led server as a macOS launchd service that starts automatically after boot.

## Overview

The `picow-led` application can be run as a macOS launchd service. This allows the web server to start automatically when the system boots, without requiring manual intervention.

## Service Configuration

The launchd service configuration file is provided as `macos-picow-led.plist`.

### Key Configuration Details

- **Label**: `com.picow-led`
- **Program**: `/usr/local/bin/picow-led server`
- **Port**: `:50836` (as specified in the original systemd service)
- **Path Prefix**: `/picow-led`
- **Log Format**: `text`
- **Database Path**: `%h/.config/picow-led/picow-led.db` (expands to user's home directory)
- **Start Behavior**: Does not run at load (`RunAtLoad=false`)

## Installation Steps

1. **Build the picow-led executable** (if not already built):
   ```bash
   cd /Users/knackwurstking/Git/picow-led
   go build -o picow-led cmd/picow-led/main.go
   ```

2. **Install the executable to a system path**:
   ```bash
   sudo cp picow-led /usr/local/bin/picow-led
   sudo chmod +x /usr/local/bin/picow-led
   ```

3. **Create the launchd service file**:
   ```bash
   # Create in user space (preferred for user services)
   mkdir -p ~/Library/LaunchAgents
   cp cmd/picow-led/services/macos-picow-led.plist ~/Library/LaunchAgents/com.picow-led.plist
   ```

4. **Load and start the service**:
   ```bash
   # Load the service
   launchctl load ~/Library/LaunchAgents/com.picow-led.plist

   # Start the service immediately (optional)
   launchctl start com.picow-led
   ```

## Service Management

After installation, you can manage the service with these commands:

- **Load the service**:
  ```bash
  launchctl load ~/Library/LaunchAgents/com.picow-led.plist
  ```

- **Unload the service**:
  ```bash
  launchctl unload ~/Library/LaunchAgents/com.picow-led.plist
  ```

- **Start the service**:
  ```bash
  launchctl start com.picow-led
  ```

- **Stop the service**:
  ```bash
  launchctl stop com.picow-led
  ```

- **Check service status**:
  ```bash
  launchctl print com.picow-led
  ```

## Log Files

The service logs are written to:
- `/var/log/picow-led.log`

## Troubleshooting

1. **Service not starting**: Check logs in `/var/log/picow-led.log`
2. **Permission issues**: Ensure the executable has proper permissions (`chmod +x`)
3. **Path issues**: Verify that `/usr/local/bin/picow-led` exists and is executable
4. **Network issues**: The service requires network connectivity to start properly

## Notes

- The service is configured to run as a user-level agent (`~/Library/LaunchAgents/`) rather than system-wide daemon (`/Library/LaunchDaemons/`)
- The service does not automatically restart on failure (KeepAlive=false)
- Network connectivity is required for the service to function properly
- The database path uses `%h` which expands to the user's home directory

## Configuration Customization

To customize the service, edit `~/Library/LaunchAgents/com.picow-led.plist`:

- Change port by modifying the `-addr` argument
- Modify database path by changing the `-database-path` argument
- Change log location by modifying `StandardOutPath` and `StandardErrorPath`
