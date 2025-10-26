import "vite-plugin-pwa/info";
import "vite-plugin-pwa/svelte";

// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        // interface Locals {}
        // interface PageData {}
        // interface PageState {}
        // interface Platform {}
    }

    interface Window {
        ws: import("ui").WS<WSMessageData>;
    }

    type PowerState = 0 | 1;

    type Device = {
        addr: string;
        name: string;
        color: number[];
        pins: number[];
        active_color: number[];
        power: PowerState;
    };

    type Devices = Device[];

    type Color = {
        id?: number;
        r: number;
        g: number;
        b: number;
    };

    type Colors = Color[];

    type WSMessageData =
        | {
              type: "devices";
              data: Devices;
          }
        | {
              type: "device";
              data: Device;
          }
        | {
              type: "colors";
              data: Colors;
          };
}

export {};
