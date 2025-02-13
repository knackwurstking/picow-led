import * as ui from "ui";

import * as alerts from "./alerts";
import * as types from "./types";
import * as ws from "./ws";

export const version = "v0.9.0";
export const store = createStore();
export const websocket = createWebSocket(store);

function createStore(): types.PicoWStore {
    const storePrefix = "picow:";
    const store: types.PicoWStore = new ui.Store(storePrefix);

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
}

function createWebSocket(store: types.PicoWStore): ws.WS {
    const websocket = new ws.WS();

    store.listen(
        "server",
        (data) => {
            console.debug("Connect WebSocket to", data);
            websocket.server = data;
        },
        true,
    );

    const statusLED = document.querySelector<HTMLElement>(`.status-led`)!;

    websocket.events.addListener("open", () => {
        console.debug("WebSocket opened");
        statusLED.setAttribute("active", "");
        statusLED.nextElementSibling!.innerHTML = "Online";

        if (store.get("firstTimeConnect") === true) {
            store.set("firstTimeConnect", false);
        }
    });

    websocket.events.addListener("close", () => {
        console.debug("WebSocket closed");
        statusLED.removeAttribute("active");
        statusLED.nextElementSibling!.innerHTML = "Offline";
    });

    websocket.events.addListener("message-error", (msg) => {
        console.error(`WebSocket "message-error":`, msg);
        alerts.add("error", msg);
    });

    return websocket;
}
