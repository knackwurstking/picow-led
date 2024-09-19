# PicoW LED Server

- [Notes](#notes)
- [TODOs](#todos)

<a id="notes"></a>

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

<a id="todos"></a>

## TODOs

- [x] Devices page
- [ ] Groups page
- [ ] Scenes page
- [x] Settings page
- [ ] Colors page - add, remove or change api colors
- [x] Need to add protocol select to settings (ssl checkbox)
- [x] Find an capacitor option to allow unsecure http connections for websockets
- [x] Adding android app icons
- [ ] Adding android app icons for dark mode "icon-dark.png"
- [x] Adding ios app icons
- [ ] Adding ios app icons for dark mode "icon-dark.png"
