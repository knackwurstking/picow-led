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

- [ ] Establish a new package `package control` with a global registry-like object (WIP).
- [ ] Upon server startup, establish connections and configurations for all devices.
- [ ] Pass the registry-like object from the `package control` to the `package services` handlers.
- [ ] After each change in services (get, update), synchronize the devices. Ensure that the handler or equivalent is passed to the service.
