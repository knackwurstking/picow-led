# PicoW LED Server

<!--toc:start-->

- [PicoW LED Server](#picow-led-server)
  - [Notes](#notes)

<!--toc:end-->

## Notes

**Routing**:

| Endpoint          | GET | POST | PUT | DELETE |
| ----------------- | :-: | :--: | :-: | :----: |
| /events/device    |  x  |      |     |        |
| /events/devices   |  x  |      |     |        |
| /api              |  x  |      |     |        |
| /api/devices      |  x  |      |     |        |
| /api/device       |  x  |  x   |  x  |   x    |
| /api/device/pins  |  x  |  x   |     |        |
| api/device/color  |  x  |  x   |     |        |
| /api/colors       |  x  |      |     |        |
| /api/colors/:name |  x  |  x   |  x  |   x    |

**API**:

```json
{
    "devices": [
        {
            "server": {
                "name": "Picow Test Device",
                "addr": "192.168.178.58:3000"
            },
            "pins": [0, 1, 2, 3],
            "color": [0, 0, 0, 0]
        }
    ],
    "colors": {
        "red": [255, 0, 0, 0],
        "green": [0, 255, 0, 0],
        "blue": [0, 0, 255, 0],
        "white": [255, 255, 255, 255]
    }
}
```
