# CHANGELOG

<!--toc:start-->

- [CHANGELOG](#changelog)
  - [v0.8.0 — [2024-10-22]](#v080-2024-10-22)
  - [v0.7.0 — [2024-10-09]](#v070-2024-10-09)
  - [v0.6.2 — [2024-09-30]](#v062-2024-09-30)
  - [v0.6.1 — [2024-09-30]](#v061-2024-09-30)

<!--toc:end-->

## v0.8.0 — [unreleased]

**General**:

- [frontend] Updated
    ["ui" to version v0.9.0](https://github.com/knackwurstking/ui)
    , using Lit web components now
- [backend] Moved "ws" package to models/

## v0.7.0 — [2024-10-09]

- [ui v0.8.0](https://github.com/knackwurstking/ui)

**General**:

- [backend] Completely rewritten, using only websockets now
- [frontend] Converted to typescript, updated to fit the new backend

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
