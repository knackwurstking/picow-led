# CHANGELOG

<!--toc:start-->

- [CHANGELOG](#changelog)
  - [v0.8.2 — [unreleased]](#v082-unreleased)
  - [v0.8.1 — [2024-10-24 Quick Fix]](#v081-2024-10-24-quick-fix)
  - [v0.8.0 — [2024-10-24]](#v080-2024-10-24)
  - [v0.7.0 — [2024-10-09]](#v070-2024-10-09)
  - [v0.6.2 — [2024-09-30]](#v062-2024-09-30)
  - [v0.6.1 — [2024-09-30]](#v061-2024-09-30)

<!--toc:end-->

## v0.8.2 — [unreleased]

**Removed**:

- Removed scrollbars from pages

## v0.8.1 — [2024-10-24 Quick Fix]

**Fixed**:

- Fixed wrong fonts path for `includeAssets` (PWA manifest), and themes added

## v0.8.0 — [2024-10-24]

**General**:

- [frontend] Updated
    ["ui" to version v1.0.0](https://github.com/knackwurstking/ui)
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
