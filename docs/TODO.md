# TODO

## Project

- [ ] Think about separating the home handler in to more handlers, one for devices and one for groups?
- [ ] Improve error handling and create a errors package

## Features

- [ ] Change group power on/off to toggle, just like devices
- [ ] OOB render group power buttons
- [ ] Find a good way to show devices overview, show power on/off, color, pins

## Performance

- **Unnecessary Database Operations**: The `services/device-controls.go` file has a method called `GetPins` that calls `control.GetPins` which seems to be a separate function. This could be optimized by caching or pre-fetching pins.

- **Caching**: There's no caching mechanism for frequently accessed data, which could lead to performance issues.
