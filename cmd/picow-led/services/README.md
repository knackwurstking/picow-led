# Services Directory

This directory contains service configuration files and documentation for running the picow-led application on different platforms.

## Contents

### macOS Launchd Service
- `macos-picow-led.plist`: Launchd service configuration for macOS
- `macos-service-setup.md`: Comprehensive documentation for setting up the macOS service

## Service Configuration

The picow-led application can be configured to run as a system service:
- **Port**: `:50836` (matches the systemd service configuration)
- **Path Prefix**: `/picow-led`
- **Log Format**: `text`
- **Database Path**: `%h/.config/picow-led/picow-led.db` (expands to user's home directory)

## Platform Support

Currently supports:
- macOS launchd services

Planned support:
- Linux systemd services
- Windows service configuration
- Docker container configuration

## Usage

The service files in this directory provide the necessary configuration to run picow-led as a background service that starts automatically after system boot.

Each platform-specific service includes appropriate configuration for:
- Process execution
- Log handling
- Resource management
- Startup behavior
