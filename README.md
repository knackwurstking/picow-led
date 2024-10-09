# PicoW LED Server

<!--toc:start-->

- [PicoW LED Server](#picow-led-server)
  - [Notes](#notes)

<!--toc:end-->

## Notes

**Endpoints**:

| Endpoint | GET | POST | PUT | DELETE |
| -------- | :-: | :--: | :-: | :----: |
| /        |  x  |      |     |        |
| /ws      |  x  |      |     |        |

**WebSocket Events**:

| Type    | Data                             |
| ------- | -------------------------------- |
| error   | `string`                         |
| devices | [`WSEvents_Device[]`](#ts-types) |
| device  | [`WSEvents_Device`](#ts-types)   |

**WebSocket Commands**:

| Command               | Data                                |
| --------------------- | ----------------------------------- |
| GET api.devices       | `null`                              |
| POST api.device       | [`WSEvents_Device`](#ts-types)      |
| PUT api.device        | [`WSEvents_Device`](#ts-types)      |
| DELETE api.device     | `{ addr: string }`                  |
| POST api.device.pins  | `{ addr: string; pins: number[] }`  |
| POST api.device.color | `{ addr: string; color: number[] }` |

<a id="ts-types"></a>

**TS Types**:

_[frontend/src/lib/websocket/types.ts](frontend/src/lib/websocket/types.ts)_

```typescript
export interface WSEvents_Device {
    server: WSEvents_DeviceServer;
    pins?: number[];
    color?: number[];
}

export interface WSEvents_DeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}
```
