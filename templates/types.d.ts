declare type PageWindow = Window & typeof globalThis & {
    ui: UI;
    api: Api;
    ws: WS;
    utils: Utils;
};

declare type UI = typeof import("ui");

declare type UIStore = import("ui").Store<{ devices: Device[] }>;

declare type Api = {
    devices: () => Promise<Device[]>;
    setDevicesColor: (color: MicroColor | undefined | null, ...devices: Device[]) => Promise<Device[]>;
};

declare type WS = {
    addr: () => string;
    isOpen: () => boolean;
    connect: () => void;
    close: () => void;
};

declare type Utils = {
    onClickPowerButton: (ev: Event & { currentTarget: HTMLButtonElement }) => Promise<void>;
    updateDeviceListItem: (item: HTMLElement, device: Device) => void;
    setupAppBarItems: (...itemNames: AppBarItemName[]) => AppBarItems;
    setOnlineIndicatorState: (state: boolean) => void;
    registerServiceWorker: () => void;
};

declare type AppBarItemName =
    | "back-button"
    | "online-indicator"
    | "title"
    | "settings-button";

declare type AppBarItems = Partial<Record<AppBarItemName, HTMLElement>>;

declare type Device = {
    server: Server;
    online: boolean;
    error: string;
    color: MicroColor;
    pins: MicroPins;
};

declare type Server = {
    addr: string;
    name: string;
};

declare type MicroColor = number[]
declare type MicroPins = number[] 
