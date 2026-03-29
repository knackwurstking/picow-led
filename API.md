# API Documentation

<!--toc:start-->

- [API Documentation](#api-documentation)
  - [Device Control Endpoints](#device-control-endpoints)
    - [Set Device Color (RGB)](#set-device-color-rgb)
    - [Set Device White Channel](#set-device-white-channel)
    - [Set Device RGBW (Color + White)](#set-device-rgbw-color-white)
    - [Set Device Brightness (Not Implemented)](#set-device-brightness-not-implemented)
    - [Set Device White2 (Not Implemented)](#set-device-white2-not-implemented)

<!--toc:end-->

## Device Control Endpoints

### Set Device Color (RGB)

```
POST /api/devices/{id}/color?color=R,G,B
```

- **Parameters**: `color` (comma-separated RGB values 0-255)
- **Returns**: `204 No Content`
- **Device Type**: RGB devices only, as long as the device type contains "RGB"
  it should be fine to use
- **Example**: `/api/devices/1/color?color=255,0,0`

### Set Device White Channel

```
POST /api/devices/{id}/white?white=VALUE
```

- **Parameters**: `white` (0-255)
- **Returns**: `204 No Content`
- **Device Type**: White channel devices only
- **Example**: `/api/devices/1/white?white=255`

### Set Device RGBW (Color + White)

```
POST /api/devices/{id}/rgbw?color=R,G,B&white=VALUE
```

- **Parameters**: `color` (RGB), `white` (0-255)
- **Returns**: `204 No Content`
- **Device Type**: RGBW devices only
- **Example**: `/api/devices/1/rgbw?color=255,0,0&white=128`

### Set Device Brightness (Not Implemented)

```
POST /api/devices/{id}/brightness
```

- **Status**: Not implemented (returns 501)

### Set Device White2 (Not Implemented)

```
POST /api/devices/{id}/white2
```

- **Status**: Not implemented (returns 501)
