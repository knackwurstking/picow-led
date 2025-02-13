import * as ui from "ui";
import * as types from "./ws/types";

export interface Component<T extends HTMLElement> {
    element: T;
}

export interface ColorStore {
    /**
     * Cache the last 10 items here or so (RGB|RGBW...)
     */
    cache: number[][];
    devices: {
        [key: string]: number[];
    };
}

export type PicoWStore = ui.Store<{
    server: types.WSServer;
    firstTimeConnect: boolean;
    color: ColorStore;
}>;
