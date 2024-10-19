import type { UIStackLayout, UIStore, UIThemeHandlerThemes } from "ui";

import type { WSEventsDevice, WSEventsServer } from "./lib/websocket";

export type PicowStackLayoutPages = "devices" | "settings";
export type PicowStackLayout = UIStackLayout<PicowStackLayoutPages>;
export type PicowStore = UIStore<PicowStoreEvents>;

export interface PicowStoreEvents {
    currentPage: PicowStackLayoutPages;
    devices: WSEventsDevice[];
    devicesColor: { [key: string]: number[] };
    server: WSEventsServer;
    currentTheme: {
        theme: UIThemeHandlerThemes;
    };
}
