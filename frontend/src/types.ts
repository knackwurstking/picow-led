/// <reference types="vite-plugin-pwa/client" />

import type { UIAppBar, UIStore, UIThemeHandlerTheme } from "ui";

export type PicowStackLayoutPage = "devices" | "settings" | "";
export type PicowStore = UIStore<PicowStoreEvents>;
export type PGAppBar = UIAppBar<"menu" | "status" | "title" | "add">;

export type WSEventsCommand = {
    "GET api.devices": null;
    "POST api.device": WSEventsDevice;
    "PUT api.device": WSEventsDevice;
    "DELETE api.device": { addr: string };
    "POST api.device.pins": { addr: string; pins: number[] };
    "POST api.device.color": { addr: string; color: number[] };
};

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

export interface PicowStoreEvents {
    currentPage: PicowStackLayoutPage;
    devices: WSEventsDevice[];
    devicesColor: { [key: string]: number[] };
    server: WSEventsServer;
    currentTheme: {
        theme: UIThemeHandlerTheme;
    };
}

export interface AppBarEvents {
    menu: Event;
    add: Event;
}

export interface WSEventsRequest {
    command: string;
    data?: string; // JSON string
}

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
