import * as globals from "../../../globals";
import * as storeUtils from "../../../store-utils";
import * as types from "../../../types";
import * as ws from "../../../ws";

import * as devicesUtils from "../utils";

const html = String.raw;
const colorStringSeparator = ", ";

export interface DeviceSetupProps {
    device: ws.WSDevice;
}

export function deviceSetup(props: DeviceSetupProps): types.Component<HTMLDialogElement> {
    const dialog = document.querySelector<HTMLDialogElement>(`dialog.device-setup`)!;

    const query = {
        title: dialog.querySelector<HTMLElement>(`.title`)!,
        form: dialog.querySelector<HTMLFormElement>(`form`)!,
        colorCache: dialog.querySelector<HTMLElement>(`.color-cache`)!,
        sliders: dialog.querySelector<HTMLUListElement>(`.sliders`)!,
        cacheButton: dialog.querySelector<HTMLButtonElement>(`button.cache`)!,
        closeButton: dialog.querySelector<HTMLButtonElement>(`button.close`)!,
    };

    query.title.innerHTML = `${props.device.server.name || props.device.server.addr}`;

    createSliders(query.sliders, props.device.pins, devicesUtils.color.get(props.device));

    // Dialog Action Buttons

    query.closeButton.onclick = () => dialog.close();

    query.cacheButton.onclick = (e) => {
        e.preventDefault();
        storeUtils.colorCache.add(getSliderValues(query.sliders));
        createColorCache(query.colorCache, query.sliders, props.device);
    };

    query.form.onsubmit = async () => {
        const color: number[] = getSliderValues(query.sliders);
        globals.store.update("color", (data) => {
            data.devices[props.device.server.addr] = color;
            return data;
        });

        await globals.websocket.request("POST api.device.color", {
            addr: props.device.server.addr,
            color,
        });
    };

    createColorCache(query.colorCache, query.sliders, props.device);

    return {
        element: dialog,
    };
}

function createSliders(
    container: HTMLUListElement,
    pins: number[] | null | undefined,
    color: number[],
): void {
    container.innerHTML = "";

    (pins || []).forEach((p, i) => {
        // create slider
        const item = document.createElement("li");

        item.className = "ui-flex-grid-row";
        item.style.setProperty("--align", "center");

        item.innerHTML = html`
            <label style="font-size: 1.2rem;" for="pin${p}">${p}</label>
            <input id="pin${p}" type="range" min="0" max="255" value="${color[i]}" />

            <label class="value" style="min-width: 3ch;" for="pin${p}"> ${color[i]} </label>
        `;

        const label = item.querySelector<HTMLElement>(`label.value`)!;
        const input = item.querySelector<HTMLInputElement>(`input`)!;
        input.oninput = () => {
            label.innerHTML = input.value;
        };

        container.appendChild(item);
    });
}

function getSliderValues(slidersContainer: HTMLUListElement): number[] {
    const color: number[] = [];
    (Array.from(slidersContainer.children) as HTMLLIElement[]).forEach((c) =>
        color.push(parseInt(c.querySelector(`input`)!.value, 10)),
    );
    return color;
}

function createColorCache(
    container: HTMLElement,
    slidersContainer: HTMLUListElement,
    device: ws.WSDevice,
) {
    container.innerHTML = "";

    storeUtils.colorCache.getAll().forEach((color) => {
        const item = document.createElement("div");

        item.style.setProperty("--align", "center");
        item.style.cursor = "pointer";
        item.style.flex = "0";

        item.className = "item ui-flex-grid ui-border ui-none-select";

        item.tabIndex = 0;
        item.role = "button";

        item.oncontextmenu = (e) => {
            e.preventDefault();
            if (confirm(`You want to delete this item: ${color.join(colorStringSeparator)}?`)) {
                storeUtils.colorCache.remove(color);
                createColorCache(container, slidersContainer, device);
            }
        };

        item.onclick = () => {
            const newColor = (Array.from(item.querySelectorAll(`.pin`)) as HTMLElement[]).map((c) =>
                parseInt(c.innerText, 10),
            );

            createSliders(
                slidersContainer,
                device.pins,
                newColor.map((c) => (isNaN(c) ? 0 : c)),
            );
        };

        item.innerHTML = html`
            <span>
                ${color.map((c, i) => html`<div class="pin pin${i}">${c}</div>`).join("")}
            </span>

            <div
                class="preview ui-border"
                style="
                    background: rgb(${color.slice(0, 3).join(", ")});
                    width: 2.5rem;
                    height: 2.5rem;
                "
            ></div>
        `;

        container.appendChild(item);
    });
}
