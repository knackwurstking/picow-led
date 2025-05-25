# PicoW LED

<!--toc:start-->
- [PicoW LED](#picow-led)
  - [Routes](#routes)
    - [/api/devices [GET]](#apidevices-get)
    - [/api/devices/:addr [GET]](#apidevicesaddr-get)
    - [/api/devices/:addr/color [GET]](#apidevicesaddrcolor-get)
    - [/api/devices/:addr/pins [GET]](#apidevicesaddrpins-get)
    - [/api/devices/:addr/name [GET]](#apidevicesaddrname-get)
    - [/api/devices/:addr/power [GET]](#apidevicesaddrpower-get)
    - [/api/devices/:addr/power [POST]](#apidevicesaddrpower-post)
    - [/api/colors [GET]](#apicolors-get)
    - [/api/colors [POST]](#apicolors-post)
    - [/api/colors [PUT]](#apicolors-put)
    - [/api/colors/:index [GET]](#apicolorsindex-get)
    - [/api/colors/:index [PUT]](#apicolorsindex-put)
    - [/api/colors/:index [DELETE]](#apicolorsindex-delete)
<!--toc:end-->

## Routes

### /api/devices [GET]

Request:

```bash
curl http://localhost:50835/api/devices
```

Response:

```json
[
    {
        "addr": "192.168.178.58:3000",
        "name": "Kitchen",
        "color": [0, 0, 0, 0],
        "pins": [0, 1, 2, 3],
        "power": 0,
    }
]
```

### /api/devices/:addr [GET]

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000
```

Response:

> States of power are represented by 0 for off and 1 for on.

```json
{
    "addr": "192.168.178.58:3000",
    "name": "Kitchen",
    "color": [0, 0, 0, 0],
    "pins": [0, 1, 2, 3],
    "power": 0,
}
```

### /api/devices/:addr/color [GET]

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/color
```

Response:

```json
[0, 0, 0, 0]
```

### /api/devices/:addr/pins [GET]

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/pins
```

Response:

```json
[0, 1, 2, 3]
```

### /api/devices/:addr/name [GET]

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/name
```

Response:

```json
"Kitchen"
```

### /api/devices/:addr/power [GET]

PowerStates: 0 | 1

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/power
```

Response:

```json
1
```

### /api/devices/:addr/power [POST]

PowerStates: 0 | 1

Request:

```bash
curl -X POST http://localhost:50835/api/devices/192.168.178.58:3000/power?state=0
```

### /api/colors [GET]

Request:

```bash
curl http://localhost:50835/api/colors
```

Response:

```json
[
    { "r": 255 "g": 255 "b": 255 },
    { "r": 255 "g": 0   "b": 0   },
    { "r": 0   "g": 255 "b": 0   },
    { "r": 0   "g": 0   "b": 255 },
]
```

### /api/colors [POST]

Request:

```bash
curl -X POST http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 255,"b": 255 }]'
```

### /api/colors [PUT]

Request:

```bash
curl -X PUT http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 0,"b": 0 }, { "r": 0,"g": 255,"b": 0 }, { "r": 0,"g": 0,"b": 255 }]'
```

### /api/colors/:index [GET]

Request:

```bash
curl http://localhost:50835/api/colors/0
```

Response:

```json
{ "r": 255 "g": 255 "b": 255 }
```

### /api/colors/:index [PUT]

Request:

```bash
curl -X PUT http://localhost:50835/api/colors/0 \
    -H "Content-Type: application/json" \
    -d '{ "r": 150, "g": 150, "b": 150 }'
```

### /api/colors/:index [DELETE]

Request:

```bash
curl -X DELETE http://localhost:50835/api/colors/0
```
