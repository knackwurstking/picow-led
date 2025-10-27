# TODO

- [ ] SQLite Database [WIP]
  - [x] Open the database, take path from flags
  - [ ] Create tables and models [WIP]
    - [x] Create a table for devices
    - [x] Create a table for device_setups
    - [ ] Create a table for duty (id, name, duty, created_at)
      - Duty contains the an rgb like values (0-255)
    - [ ] Create a table for groups (id, name, setup, created_at)
      - Setup contains a list of device_id and duty data
    - [ ] Create a table for scenes (id, name, device_ids, created_at)

- [ ] Create the UI and use the echo router
- [ ] Using the assets system from the previous project (pg-press) or just add this to the ui lib repo [WIP]
