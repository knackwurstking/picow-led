// @ts-ignore
import svgPower from "ui/svg/power";

import * as globals from "../../../globals";
import * as types from "../../../types";
import * as ws from "../../../ws";

import * as dialogs from "../dialogs";
import * as devicesUtils from "../utils";

const html = String.raw;
const colorStringSeparator = ", ";

// TODO: Mark offline device somehow
export function deviceListItem(device: ws.WSDevice): types.Component<HTMLLIElement> {
    const powerButtonColor = devicesUtils.getPowerButtonColor(device.color);
    const item = document.createElement("li");

    item.setAttribute("data-addr", `${device.server.addr}`);
    item.setAttribute("data-current-color", JSON.stringify(device.color));

    item.className = "device";

    item.innerHTML = html`
        <fieldset>
            <legend>${device.server.name || device.server.addr}</legend>

            <div class="ui-flex-grid-row" style="--align: center">
                <div class="ui-flex-grid-item">
                    <code>
                        ${device.color?.map((c) => c.toString()).join(colorStringSeparator) || ""}
                    </code>
                </div>

                <div class="ui-flex-grid-item" style="--flex: 0">
                    <button
                        class="power"
                        style="
                            overflow: visible;
                            height: 3rem;
                            width: 3rem;
                            transition: color 0.25s linear;
                            color: ${powerButtonColor};
                        "
                        variant="ghost"
                        icon
                    >
                        <div
                            class="background"
                            style="
                                position: absolute;
                                top: calc(var(--ui-spacing) / 2 - 1px);
                                right: calc(var(--ui-spacing) / 2 - 1px);
                                bottom: calc(var(--ui-spacing) / 2 - 1px);
                                left: calc(var(--ui-spacing) / 2 - 1px);
                                border-radius: 50%;
                                background-color: ${powerButtonColor};
                                filter: blur(2.5px);
                                overflow: visible;
                            "
                        ></div>
                        ${svgPower}
                    </button>
                </div>

                <div class="ui-flex-grid-item" style="--flex: 0">
                    <button class="options" variant="ghost" color="primary" icon>
                        <div class="bi bi-three-dots-vertical"></div>
                    </button>
                </div>
            </div>
        </fieldset>
    `;

    initPowerButton(item, device);
    initOptionsButton(item, device);

    devicesUtils.color.set(item.getAttribute("data-addr")!, device.color);

    return {
        element: item,
    };
}

// TODO: Mark offline device somehow
export function updateDeviceListItem(item: HTMLLIElement | HTMLElement, device: ws.WSDevice): void {
    if (item.getAttribute("data-addr") !== device.server.addr) {
        return;
    }

    item.setAttribute("data-current-color", JSON.stringify(device.color));

    // Update the color preview code block
    const code = item.querySelector(`code`)!;
    code.innerText = device.color?.map((c) => c.toString()).join(colorStringSeparator) || "";

    // Update the power toggle button color
    const powerButton = item.querySelector<HTMLElement>(`button.power`)!;
    powerButton.style.color = devicesUtils.getPowerButtonColor(device.color);

    powerButton.querySelector<HTMLElement>(`.background`)!.style.background =
        devicesUtils.getPowerButtonColor(device.color);

    devicesUtils.color.set(item.getAttribute("data-addr")!, device.color);
}

function initPowerButton(item: HTMLLIElement, device: ws.WSDevice) {
    const powerButton = item.querySelector<HTMLElement>(`button.power`)!;
    powerButton.onclick = async () => {
        powerButton.style.color = "yellow"; // Work in progress

        const addr: string = item.getAttribute("data-addr")!;
        let color: number[] = JSON.parse(item.getAttribute("data-current-color")!);

        if (Math.max(...color) > 0) {
            color = color.map(() => 0); // Turn off
        } else {
            color = devicesUtils.color.get(device);
        }

        await globals.websocket.request("POST api.device.color", { addr, color });
    };
}

function initOptionsButton(item: HTMLLIElement, device: ws.WSDevice) {
    item.querySelector<HTMLElement>(`button.options`)!.onclick = async () => {
        const dialog = dialogs.deviceSetup({ device });
        dialog.element.showModal();
    };
}
