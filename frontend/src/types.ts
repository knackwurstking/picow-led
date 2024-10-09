import type { UIStackLayout, UIStore } from "ui";

import type { WSEvents_Device, WSEvents_Server } from "./lib/websocket";
import type { UIThemeHandler_Theme } from "ui/src/ui-theme-handler/ui-theme-handler";

export type PicowStackLayout = UIStackLayout<PicowStackLayout_Pages>;

export type PicowStackLayout_Pages = null | "devices" | "settings";

export type PicowStore = UIStore<PicowStore_Events>;

export interface DevicesColor {
    [key: string]: number[];
}

export interface PicowStore_Events {
    currentPage: PicowStackLayout_Pages;
    devices: WSEvents_Device[];
    devicesColor: DevicesColor;
    server: WSEvents_Server;
    currentTheme: {
        theme: UIThemeHandler_Theme;
    };
}
