# TODO

- [ ] SQLite Database [WIP]
  - [x] Open the database, take path from flags
  - [ ] Create tables and models [WIP]
    - [x] Create a table for devices (id, addr, name, pins, created_at)
    - [ ] Create a table for pins (id, device_id, name, pins, created_at) [WIP]
      - Pins contains a list of number (uint8)
    - [ ] Create a table for colors (id, name, duty, created_at)
      - Duty contains the an rgb like values (0-255)
    - [ ] Create a table for groups (id, name, setup, created_at)
      - Setup contains a list of device_id and duty data
    - [ ] Create a table for scenes (id, name, device_ids, created_at)

- [ ] Create the UI and use the echo router
- [ ] Using the assets system from the previous project (pg-press) or just add this to the ui lib repo [WIP]
