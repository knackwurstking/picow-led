export interface WSServer {
    ssl: boolean;
    host: string;
    port: string;
}

export interface WSDeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}

export interface WSDevice {
    server: WSDeviceServer;
    pins?: number[];
    color?: number[];
}

export type WSCommand = {
    "GET api.devices": null;
    "POST api.device": WSDevice;
    "PUT api.device": WSDevice;
    "DELETE api.device": { addr: string };
    "POST api.device.pins": { addr: string; pins: number[] };
    "POST api.device.color": { addr: string; color: number[] };
};

export interface WSRequest {
    command: string;
    data?: string; // JSON string
}

export type WSResponse =
    | {
          data: string;
          type: "error";
      }
    | {
          data: WSDevice[];
          type: "devices";
      }
    | {
          data: WSDevice;
          type: "device";
      };
