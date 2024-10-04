export interface WSServer {
    ssl: boolean;
    host: string;
    port: string;
}

export type WSEvents_Command = {
    "GET api.devices": {
        request: null;
        response: WSEvents_Device[];
    };
    "POST api.device.color": {
        request: {
            addr: string;
            color: number[];
        };
        response: null;
    };
};

export interface WSEvents_DeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}

export interface WSEvents_Device {
    server: WSEvents_DeviceServer;
    pins?: number[];
    color?: number[];
}

export interface WSEvents_Request {
    command: string;
    data: string; // NOTE: JSON string
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
