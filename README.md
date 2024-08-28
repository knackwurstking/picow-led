# PicoW LED Server

- [Notes](#notes)

## <a id="notes"></a>Notes

- Using the "echo" server package
- Logging via "slog" using "lmittmann/tint"

**Routing**:

| Endpoint                | Description                                       |
| ----------------------- | ------------------------------------------------- |
| /api                    | ???                                               |
| /api/config             | GET/PUT - Server configuration                    |
| /api/config/range       | GET/POST - Global range configuration             |
| /api/config/devices     | GET/POST/PUT/DELETE - Whole device configurations |
| /api/config/devices/:id | PUT - Modify device configurations                |
| /api/devices            | POST device control commands (Multiple devices)   |
| /api/devices/:id        | POST device control commands (Single device)      |
| /api/colors             | GET - Colors                                      |
| /api/colors/:name       | GET/POST - Color                                  |

**API**:

```json
{
  "config": {
    "range": {
      "min": 0,
      "max": 100
    },
    "devices": [
      {
        "id": 0,
        "pins": [0, 1, 2, 3],
        "range": {
          "min": 0,
          "max": 100
        }
      },
      {
        "id": 1,
        "pins": [0, 1, 2, 3],
        "range": {
          "min": 0,
          "max": 100
        }
      }
    ]
  },
  "devices": [
    {
      "id": 0,
      "duty": [0, 0, 0, 0]
    },
    {
      "id": 1,
      "duty": [0, 0, 0, 0]
    }
  ],
  "colors": {
    "red": [100, 0, 0],
    "green": [0, 100, 0],
    "blue": [0, 0, 100]
  }
}
```
