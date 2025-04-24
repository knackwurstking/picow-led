import * as ui from "ui";

import * as ws from "../ws";

export type PicoWStore = ui.Store<{
    server: ws.WSServer;
    firstTimeConnect: boolean;
    color: ColorStore;
}>;

export interface ColorStore {
    /**
     * Cache the last 10 items here or so (RGB|RGBW...)
     */
    cache: number[][];
    devices: {
        [key: string]: number[];
    };
}

export const obj: PicoWStore = (() => {
    const storePrefix = "picow:";
    const store: PicoWStore = new ui.Store(storePrefix);

    store.set(
        "server",
        {
            ssl: !!location.protocol.match(/(https)/),
            host: location.hostname,
            port: location.port,
        },
        true,
    );

    store.set("firstTimeConnect", true, true);
    store.set("color", { cache: [], devices: {} }, true);

    return store;
})();

const colorStringSeparator = ",";

export const colorCache = {
    add(color: number[]): void {
        obj.update("color", (data) => {
            const newColorString = color.join(colorStringSeparator);
            if (data.cache.findIndex((c) => c.join(colorStringSeparator) === newColorString) >= 0) {
                return data;
            }

            data.cache.unshift(color);
            if (data.cache.length > 20) {
                data.cache = data.cache.slice(0, 20);
            }

            return data;
        });
    },

    getAll() {
        return obj.get("color")?.cache || [];
    },

    remove(color: number[]) {
        obj.update("color", (data) => {
            const colorString = color.join(colorStringSeparator);
            data.cache = data.cache.filter((c) => c.join(colorStringSeparator) !== colorString);
            return data;
        });
    },
};
