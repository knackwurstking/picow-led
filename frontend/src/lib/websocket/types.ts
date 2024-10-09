export type WSEvents_Command = {
    "GET api.devices": null;
    "POST api.device": WSEvents_Device;
    "PUT api.device": WSEvents_Device;
    "DELETE api.device": { addr: string };
    "POST api.device.pins": { addr: string; pins: number[] };
    "POST api.device.color": { addr: string; color: number[] };
};

export interface WSEvents_Request {
    command: string;
    data: string; // JSON string
}

export type WSEvents_Response =
    | {
          data: string;
          type: "error";
      }
    | {
          data: WSEvents_Device[];
          type: "devices";
      }
    | {
          data: WSEvents_Device;
          type: "device";
      };

export interface WSEvents_Server {
    ssl: boolean;
    host: string;
    port: string;
}

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
