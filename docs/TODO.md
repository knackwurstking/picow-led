# TODO

## Project

- [ ] Think about separating the home handler in to more handlers, one for devices and one for groups?
- [ ] Improve error handling and create a errors package

## Features

- [ ] Change group power on/off to toggle, just like devices
- [ ] OOB render group power buttons
- [ ] Find a good way to show devices overview, show power on/off, color, pins

## Performance

Unnecessary Database Operations: The `services/device-controls.go` file has a method called `GetPins` that calls `control.GetPins` which seems to be a separate function. This could be optimized by caching or pre-fetching pins.

- [x] Improve the device controls service, prefetch pins and cache it (no database needed)
- [x] Also do this for the current device color
