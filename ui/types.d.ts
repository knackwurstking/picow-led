export declare global {
    interface Window {
        ui: UI;
        store: Store;
        api: import("./lib/api").Api;
        utils: Utils;
        ws: WS;
    }

    declare type UI = typeof import("ui");

    declare type Store = {
        obj: UIStore;
        device: (addr: string) => Device | null;
        currentDeviceColor: (addr: string) => Color | null;
    };

    declare type UIStore = import("ui").Store<{
        devices: Device[];
        colors: Colors;
        currentDeviceColors: Record<string, Color>;
    }>;

    declare type Utils = {
        setupAppBarItems: (...itemNames: AppBarItemName[]) => AppBarItems;
        setOnlineIndicatorState: (state: boolean) => void;
        registerServiceWorker: () => void;
    };

    declare type AppBarItemName =
        | "online-indicator"
        | "title"
        | "settings-button";

    declare type AppBarItems = Partial<Record<AppBarItemName, HTMLElement>>;

    declare type Device = {
        server: Server;
        online: boolean;
        error: string;
        color: Color | null;
        pins: Pins | null;
    };

    declare type Server = {
        addr: string;
        name: string;
    };

    declare type Color = number[];
    declare type Pins = number[];

    declare type Colors = Color[];

    declare type WS = {
        events: import("ui").Events<{
            open: undefined;
            device: Device;
            colors: Colors;
        }>;
        socket: WebSocket | null;
        isOpen: () => boolean;
        connect: () => Promise<void>;
        close: () => void;
    };

    declare type WSMessageData =
        | {
              type: "device";
              data: Device;
          }
        | {
              type: "colors";
              data: Colors;
          };
}
