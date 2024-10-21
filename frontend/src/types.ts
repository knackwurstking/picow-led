import type { UIStore, UIThemeHandlerTheme } from "ui";

import type { WSEventsDevice, WSEventsServer } from "./lib/websocket";

export type PicowStackLayoutPage = "devices" | "settings" | "";

export type PicowStore = UIStore<PicowStoreEvents>;

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
