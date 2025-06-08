# PicoW LED Server

## v0.12.0 - 2025-06-08

- A complete rewrite with a separate frontend using sveltekit [picow-led-frontend v0.1.0](https://github.com/knackwurstking/picow-led-frontend/tree/v0.1.0)

## v0.11.6 - unreleased

- removed pageshow and pagehide event handlers
- removed /ws and /api from the service worker blacklist
- removed route /api/ping

## v0.11.5 - 2025-09-20

- upgrade ui to v4.2.1, ws fix

## v0.11.4 - 2025-09-19

- ui: restructor scripts and public/js
- ui: Use page show/hide events for scripts

## v0.11.3 - 2025-09-18

- ui: Add notifications for error messages

## v0.11.2 - 2025-09-17

- ui: Add /ws to the service-worker blacklist

## v0.11.1 - 2025-09-15

**Fixed**:

- ui: missing online indicator (offline mode)
- ui: no devices list items rendered (offline mode)

## v0.11.0 - 2025-09-14

> Depends on [picow-led-microcontroller v1.0.0](https://github.com/knackwurstking/picow-led-microcontroller#v1.0.0)
