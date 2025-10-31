# TODO

## Assets System

- [x] Utilize the assets system from the previous project (pg-press) or integrate it into the UI library repository.

## SQLite Database

- [x] Open the database using a flag-specified path.
- [x] Establish tables and models:
  - Devices
  - Device Setups
  - Colors
  - Groups
  - Device Control (WIP)
- [x] Service DeviceControl needs to be handled by the Device service, if a device got deleted, also remove the device from the DeviceControl table.

## Router, UI, Handlers

### Layout

- [x] Incorporate icons into the layout.
- [x] Include the manifest JSON file.

### Home Section: Devices

- [x] Implement the edit dialog.
- [x] Add a delete button to the `DialogEditDevice` component and update or add the corresponding handler.
- [x] Render the devices list.
- [x] Create a new device dialog (`DialogNewDevice`).
- [x] Create an edit device dialog (`DialogEditDevice`).
- [ ] Implement power toggle functionality (on/off).

### Home Section: Groups

- [ ] Create an edit dialog for groups.
- [ ] Render the groups list.
- [ ] Develop a new group dialog (`DialogNewGroup`).
- [ ] Build an edit group dialog (`DialogEditGroup`).

## Device Control Handling

- [x] Establish a new package `package control`.
- [ ] Refactor the service handlers to use control, ensuring that all operations are handled through these services.
  - [ ] Service: DeviceSetups, when changing pins
  - [ ] New DeviceControl (device_control) Service for color, version, temp and disk-usage
- [x] Improve the doc comment for the NewRequest function
- [x] Add missing info commands (info get: "temp", "disk-usage", "version")
