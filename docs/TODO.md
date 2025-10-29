# TODO

- [x] Using the assets system from the previous project (pg-press) or just add
      this to the ui lib repo

## SQLite Database

- [x] Open the database, take path from flags

Create tables and models

- [x] Create a table for devices
- [x] Create a table for device_setups
- [x] Create a table for colors
- [x] Create a table for groups

## Router, UI, Handlers

- [x] Layout: add icons for the layout
- [x] Layout: add the manifest json file

Home: Section Devices

- [x] Edit dialog
- [x] Add a delete button to the DialogEditDevice component and update or add
      the handler for that
- [x] Render devices list
- [x] Dialog: New device dialog
- [x] Dialog: Edit device dialog
- [ ] Power toggle (on/off)

Home: Section Groups

- [ ] Edit dialog
- [ ] Render groups list
- [ ] Dialog: New group dialog
- [ ] Dialog: Edit group dialog

## Handling device control

- [ ] New package `package control` with a global registry like object [WIP]
- [ ] After server started, connect and setup all devices
- [ ] Pass the registry like object from the `package control` to the
      `package services` handlers
- [ ] After each change in services (get, update), also update devices,
      need to pass the handler (or whatever) to the service
