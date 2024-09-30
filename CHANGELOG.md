# CHANGELOG

## v0.6.2 — 2024-09-30

**Fixed**:

- [frontend] Fix missing setup-device-dialog content
- [frontend] Fix missing setup-device-dialog action buttons if
    `allowDeletion` was set to false

## v0.6.1 — 2024-09-30

**Added**:

- [backend] Added `events` package ("pkg/event")

**Changed**:

- [frontend] Changed "ui" version in use to v0.7.1

**General**:

- [backend] Moved `endpoints` package to "internal/"
- [backend] Use `events.Events` for API changes for emitting WebSocket events
