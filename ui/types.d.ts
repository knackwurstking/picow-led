export declare type PageWindow = Window & typeof globalThis & {
    ui: UI;
    api: Api;
    ws: WS;
    utils: Utils;
};

export declare type UI = typeof import("ui");

export declare type UIStore = import("ui").Store<{ devices: Device[] }>;

export declare type Api = {
    devices: () => Promise<Device[]>;
    setDevicesColor: (color: Color | undefined | null, ...devices: Device[]) => Promise<Device[]>;
    color: () => Promise<ColorCache>
};

export declare type WS = {
    addr: () => string;
    isOpen: () => boolean;
    connect: () => void;
    close: () => void;
};

export declare type Utils = {
    setupAppBarItems: (...itemNames: AppBarItemName[]) => AppBarItems;
    setOnlineIndicatorState: (state: boolean) => void;
    registerServiceWorker: () => void;
};

export declare type AppBarItemName =
    | "back-button"
    | "online-indicator"
    | "title"
    | "settings-button";

export declare type AppBarItems = Partial<Record<AppBarItemName, HTMLElement>>;

export declare type Device = {
    server: Server;
    online: boolean;
    error: string;
    color: Color;
    pins: Pins;
};

export declare type Server = {
    addr: string;
    name: string;
};

export declare type Color = number[]
export declare type Pins = number[]

export declare type ColorCache = Record<string, Color>
