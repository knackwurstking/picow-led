# PicoW LED Server

- [Notes](#notes)

## <a id="notes"></a>Notes

- Using the "echo" server package
- Logging via "slog" using "lmittmann/tint"

**Routing**:

| Endpoint                | Description                                       |
| ----------------------- | ------------------------------------------------- |
| /api                    | ???                                               |
| /api/config             | GET the whole server configuration                |
| /api/config/devices     | GET/POST/PUT/DELETE - Whole device configurations |
| /api/config/devices/:id | PUT - Modify device configurations                |
| /api/devices            | POST device control commands (Multiple devices)   |
| /api/devices/:id        | POST device control commands (Single device)      |
