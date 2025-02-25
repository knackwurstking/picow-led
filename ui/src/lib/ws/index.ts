export * from "./base";
export * from "./types";
export * from "./ws";

import * as alerts from "../../alerts";
import * as store from "../store";
import * as ws from "./ws";

export const socket: ws.WS = (() => {
    const websocket = new ws.WS();

    store.obj.listen(
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

        if (store.obj.get("firstTimeConnect") === true) {
            store.obj.set("firstTimeConnect", false);
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
})();
