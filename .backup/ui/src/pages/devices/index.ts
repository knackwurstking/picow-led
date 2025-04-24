import * as ui from "ui";

import * as ws from "../../lib/ws";

import * as deviceListItem from "./device-list-item";

let cleanUp: ui.CleanUpFunction[] = [];

export async function onMount() {
    console.debug("[devices] Mount devices template/page");

    const routerTarget = document.querySelector(`.router-target`)!;
    const list = routerTarget.querySelector(`#devices-list`)!;

    cleanUp.push(
        ws.socket.events.addListener("message-devices", (devices) => {
            console.debug(`[devices] ws event "message-devices"`, devices);
            list.innerHTML = "";
            devices.forEach((d) => list.appendChild(deviceListItem.create(d)));
        }),

        ws.socket.events.addListener("message-device", (device) => {
            console.debug(`[devices] ws event "message-devices"`, device);

            for (const child of list.children) {
                deviceListItem.update(child as HTMLElement, device);
            }
        }),

        ws.socket.events.addListener("open", () => {
            console.debug("[devices] Request device list");
            ws.socket.request("GET api.devices");
        }),
    );

    ws.socket.request("GET api.devices");

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
