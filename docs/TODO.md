# TODO

## Code Quality

- [x] See if we can improve the handler and components structure

## FIX

- [x] Log: `2025-11-05 03:39:47 ERR services/utils-resolve.go:15 Failed to get current color for device id=1 error="failed to get device control: not found: sql: no rows in result set"`

## Systemd Service File and Server Setup

- [x] Create a service file for systemd
- [x] Deploy the application to the server

## Server Logs

- [x] Improve all server logs
  - services package

## Assets System

- [x] Use the assets system from the previous project (pg-press) or incorporate it into the UI library repository.

## Models

- [x] Rename the model type files for services
- [x] Ensure all JSON struct tags are present

## SQLite Database

- [x] Configure the database to open using a path specified by a flag.
- [x] Define and set up tables and models:
  - Devices
  - ~Device Setups~
  - Colors
  - Groups
  - Device Control
- [x] Handle the DeviceControl service: If a device is deleted, remove it from the DeviceControl table as well.

## Router, UI, Handlers

### Layout Enhancements

- [x] Integrate icons into the layout design.
- [x] Incorporate the manifest JSON file into the project.

### Home Section: Devices

- [x] Implement an edit dialog for devices.
- [x] Add a delete button to the `DialogEditDevice` component and add or update the corresponding handler.
- [x] Display the list of devices on the home screen.
- [x] Develop a dialog for creating a new device (`DialogNewDevice`).
- [x] Enhance the `DialogEditDevice` for editing devices.
- [x] Implement toggle functionality for power control (on/off).

### Home Section: Groups

- [ ] Build an edit dialog specifically for managing groups.
- [ ] Display a list of existing groups.
- [ ] Create a new group dialog (`DialogNewGroup`).
- [ ] Construct an edit dialog for managing group settings (`DialogEditGroup`).

## Device Control Handling

- [x] Create a new package called `control`.
- Refactor service handlers to utilize the control package, ensuring all operations go through these services:
  - [x] ~The DeviceSetups service should update the picow device after each pin change.~
  - [x] Implement a method to retrieve the current color.
  - [x] Add functionality to retrieve version information.
  - [x] Create a method for obtaining disk usage details.
  - [x] Implement a method to check temperature.
- [x] Enhance the documentation comment for the `NewRequest` function.
- [x] Add missing `info` commands for querying device properties like `temp`, `disk-usage`, and `version`.
