# PicoW LED

<!--toc:start-->
- [PicoW LED](#picow-led)
  - [Routes](#routes)
    - [GET _/api/devices_](#get-apidevices)
    - [GET _/api/devices/:addr_](#get-apidevicesaddr)
    - [GET _/api/devices/:addr/color_](#get-apidevicesaddrcolor)
    - [GET _/api/devices/:addr/pins_](#get-apidevicesaddrpins)
    - [GET _/api/devices/:addr/name_](#get-apidevicesaddrname)
    - [GET _/api/devices/:addr/power_](#get-apidevicesaddrpower)
    - [POST _/api/devices/:addr/power_](#post-apidevicesaddrpower)
    - [GET _/api/colors_](#get-apicolors)
    - [POST _/api/colors_](#post-apicolors)
    - [PUT _/api/colors_](#put-apicolors)
    - [GET _/api/colors/:index_](#get-apicolorsindex)
    - [PUT _/api/colors/:index_](#put-apicolorsindex)
    - [DELETE _/api/colors/:index_](#delete-apicolorsindex)
<!--toc:end-->

## Routes

### GET _/api/devices_

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

### GET _/api/devices/:addr_

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

### GET _/api/devices/:addr/color_

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/color
```

Response:

```json
[0, 0, 0, 0]
```

### GET _/api/devices/:addr/pins_

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/pins
```

Response:

```json
[0, 1, 2, 3]
```

### GET _/api/devices/:addr/name_

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/name
```

Response:

```json
"Kitchen"
```

### GET _/api/devices/:addr/power_

PowerStates: 0 | 1

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/power
```

Response:

```json
1
```

### POST _/api/devices/:addr/power_

PowerStates: 0 | 1

Request:

```bash
curl -X POST http://localhost:50835/api/devices/192.168.178.58:3000/power?state=0
```

### GET _/api/colors_

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

### POST _/api/colors_

Request:

```bash
curl -X POST http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 255,"b": 255 }]'
```

### PUT _/api/colors_

Request:

```bash
curl -X PUT http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 0,"b": 0 }, { "r": 0,"g": 255,"b": 0 }, { "r": 0,"g": 0,"b": 255 }]'
```

### GET _/api/colors/:index_

Request:

```bash
curl http://localhost:50835/api/colors/0
```

Response:

```json
{ "r": 255 "g": 255 "b": 255 }
```

### PUT _/api/colors/:index_

Request:

```bash
curl -X PUT http://localhost:50835/api/colors/0 \
    -H "Content-Type: application/json" \
    -d '{ "r": 150, "g": 150, "b": 150 }'
```

### DELETE _/api/colors/:index_

Request:

```bash
curl -X DELETE http://localhost:50835/api/colors/0
```
