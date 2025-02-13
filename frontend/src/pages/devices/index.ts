import * as ui from "ui";

import * as globals from "../../globals";
import * as create from "./create";

let cleanUp: ui.CleanUpFunction[] = [];

export async function onMount() {
    console.debug("[devices] Mount devices template/page");

    const routerTarget = document.querySelector(`.router-target`)!;
    const list = routerTarget.querySelector(`#devices-list`)!;

    cleanUp.push(
        globals.websocket.events.addListener("message-devices", (devices) => {
            console.debug(`[devices] ws event "message-devices"`, devices);
            list.innerHTML = "";
            devices.forEach((d) => list.appendChild(create.deviceListItem(d).element));
        }),

        globals.websocket.events.addListener("message-device", (device) => {
            console.debug(`[devices] ws event "message-devices"`, device);

            for (const child of list.children) {
                create.updateDeviceListItem(child as HTMLElement, device);
            }
        }),

        globals.websocket.events.addListener("open", () => {
            console.debug("[devices] Request device list");
            globals.websocket.request("GET api.devices");
        }),
    );

    globals.websocket.request("GET api.devices");

    const settingsButton = document.querySelector<HTMLElement>(`.ui-app-bar button#goToSettings`)!;
    settingsButton.style.display = "block";
    settingsButton.onclick = () => {
        location.hash = "#settings";
    };
}

export async function onDestroy() {
    console.debug("[devices] Destroy devices template/page");

    cleanUp.forEach((fn) => fn());
    cleanUp = [];

    const settingsButton = document.querySelector<HTMLElement>(`.ui-app-bar button#goToSettings`)!;
    settingsButton.style.display = "none";
}
