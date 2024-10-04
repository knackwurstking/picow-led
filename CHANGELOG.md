# CHANGELOG

<!--toc:start-->

- [CHANGELOG](#changelog)
  - [v0.7.0 — [unreleased]](#v070-unreleased)
  - [v0.6.2 — [2024-09-30]](#v062-2024-09-30)
  - [v0.6.1 — [2024-09-30]](#v061-2024-09-30)

<!--toc:end-->

## v0.7.0 — [unreleased]

**General**:

- [backend] Completely rewritten, using only websockets now
- [frontend] Converted to typescript

## v0.6.2 — [2024-09-30]

**Fixed**:

- [frontend] Fix missing setup-device-dialog content
- [frontend] Fix missing setup-device-dialog action buttons if
    `allowDeletion` was set to false

## v0.6.1 — [2024-09-30]

**Added**:

- [backend] Added `events` package ("pkg/event")

**Changed**:

- [frontend] Changed "ui" version in use to v0.7.1

**General**:

- [backend] Moved `endpoints` package to "internal/"
- [backend] Use `events.Events` for API changes for emitting WebSocket events
