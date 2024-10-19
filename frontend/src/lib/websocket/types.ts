export type WSEventsCommand = {
    "GET api.devices": null;
    "POST api.device": WSEventsDevice;
    "PUT api.device": WSEventsDevice;
    "DELETE api.device": { addr: string };
    "POST api.device.pins": { addr: string; pins: number[] };
    "POST api.device.color": { addr: string; color: number[] };
};

export interface WSEventsRequest {
    command: string;
    data: string; // JSON string
}

export type WSEventsResponse =
    | {
          data: string;
          type: "error";
      }
    | {
          data: WSEventsDevice[];
          type: "devices";
      }
    | {
          data: WSEventsDevice;
          type: "device";
      };

export interface WSEventsServer {
    ssl: boolean;
    host: string;
    port: string;
}

export interface WSEventsDevice {
    server: WSEventsDeviceServer;
    pins?: number[];
    color?: number[];
}

export interface WSEventsDeviceServer {
    name?: string;
    addr: string;
    online?: boolean;
}
