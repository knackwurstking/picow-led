# TODO

## Project

- [ ] Think about separating the home handler into more handlers, one for devices and one for groups?
- [ ] Improve error handling and create a errors package

## Fixes

- [x] Remove the color cache, only keep the pins cache, also check the expiration validation [WIP]
- [x] Makefile: Database path for the macos service, Need to use the correct path: ~/Library/Application\ Support/picow-led.

## Performance

Unnecessary Database Operations: The `services/device-controls.go` file has a method called `GetPins` that calls `control.GetPins` which seems to be a separate function. This could be optimized by caching or pre-fetching pins.

- [x] Improve the device controls service, prefetch pins and cache it (no database needed)
- [x] Also do this for the current device color
