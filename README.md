# PicoW LED

<!--toc:start-->
- [PicoW LED](#picow-led)
  - [Routes](#routes)
    - [_/api/devices_ **GET**](#apidevices-get)
    - [_/api/devices/:addr_ **GET**](#apidevicesaddr-get)
    - [_/api/devices/:addr/name_ **GET**](#apidevicesaddrname-get)
    - [_/api/devices/:addr/color_ **GET**](#apidevicesaddrcolor-get)
    - [_/api/devices/:addr/pins_ **GET**](#apidevicesaddrpins-get)
    - [_/api/devices/:addr/active_color_ **GET**](#apidevicesaddractivecolor-get)
    - [_/api/devices/:addr/power_ **GET**](#apidevicesaddrpower-get)
    - [_/api/devices/:addr/power_ **POST**](#apidevicesaddrpower-post)
    - [_/api/colors_ **GET**](#apicolors-get)
    - [_/api/colors_ **POST**](#apicolors-post)
    - [_/api/colors_ **PUT**](#apicolors-put)
    - [_/api/colors/:index_ **GET**](#apicolorsindex-get)
    - [_/api/colors/:index_ **PUT**](#apicolorsindex-put)
    - [_/api/colors/:index_ **DELETE**](#apicolorsindex-delete)
<!--toc:end-->

## Routes

### _/api/devices_ **GET**

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
        "active_color": [255, 255, 255, 255],
        "power": 0,
        
    }
]
```

### _/api/devices/:addr_ **GET**

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
    "active_color": [255, 255, 255, 255],
    "power": 0,
}
```

### _/api/devices/:addr/name_ **GET**

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/name
```

Response:

```json
"Kitchen"
```


### _/api/devices/:addr/color_ **GET**

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/color
```

Response:

```json
[0, 0, 0, 0]
```

### _/api/devices/:addr/pins_ **GET**

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/pins
```

Response:

```json
[0, 1, 2, 3]
```



### _/api/devices/:addr/active_color_ **GET**

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/active_color
```

Response:

```json
[255, 255, 255, 255]
```

### _/api/devices/:addr/power_ **GET**

PowerStates: 0 | 1

Request:

```bash
curl http://localhost:50835/api/devices/192.168.178.58:3000/power
```

Response:

```json
1
```

### _/api/devices/:addr/power_ **POST**

PowerStates: 0 | 1

Request:

```bash
curl -X POST http://localhost:50835/api/devices/192.168.178.58:3000/power?state=0
```

### _/api/colors_ **GET**

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

### _/api/colors_ **POST**

Request:

```bash
curl -X POST http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 255,"b": 255 }]'
```

### _/api/colors_ **PUT**

Request:

```bash
curl -X PUT http://localhost:50835/api/colors \
    -H "Content-Type: application/json" \
    -d '[{ "r": 255,"g": 0,"b": 0 }, { "r": 0,"g": 255,"b": 0 }, { "r": 0,"g": 0,"b": 255 }]'
```

### _/api/colors/:index_ **GET**

Request:

```bash
curl http://localhost:50835/api/colors/0
```

Response:

```json
{ "r": 255 "g": 255 "b": 255 }
```

### _/api/colors/:index_ **PUT**

Request:

```bash
curl -X PUT http://localhost:50835/api/colors/0 \
    -H "Content-Type: application/json" \
    -d '{ "r": 150, "g": 150, "b": 150 }'
```

### _/api/colors/:index_ **DELETE**

Request:

```bash
curl -X DELETE http://localhost:50835/api/colors/0
```
