export declare global {
    interface Window {
        ui: UI;
        store: Store;
        ws: WS;
    }
}

export declare type UI = typeof import("ui");

export declare type Store = {
    obj: UIStore;
    device: (addr: string) => Device | null;
    currentColor: (addr: string) => Color | null;
};

export declare type UIStore = import("ui").Store<{
    devices: Device[];
    color: {
        api: ColorCache; // TODO: Not handled for now (api)
        current: Record<string, Color>;
    };
}>;

export declare type Api = {
    devices: () => Promise<Device[]>;
    setDevicesColor: (
        color: Color | undefined | null,
        ...devices: Device[]
    ) => Promise<Device[]>;
    colors: () => Promise<ColorCache>;
    color: (index: number) => Promise<Color>;
    setColor: (index: number, color: Color) => Promise<void>;
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

export declare type Color = number[];
export declare type Pins = number[];

export declare type ColorCache = Color[];
