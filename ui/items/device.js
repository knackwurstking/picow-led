/** @type {import("../types.d.ts").PageWindow} */
// @ts-ignore
const w = window;

/**
 * @param {import("../types.d.ts").Device} device
 * @returns {HTMLElement}
 */
export function createDeviceItem(device) {
    /** @type {HTMLTemplateElement} */
    const template = document.querySelector(
        `template[name="device-list-item"]`,
    );

    /** @type {HTMLElement} */
    const item = template.content
        .cloneNode(true)
        // @ts-expect-error
        .querySelector(".device-list-item");

    updateDeviceItem(item, device);

    return item;
}

/**
 * @param {HTMLElement} item
 * @param {import("../types.d.ts").Device} device
 * @returns {void}
 */
export function updateDeviceItem(item, device) {
    item.setAttribute("data-addr", device.server.addr);

    /** @type {HTMLElement} */
    const title = item.querySelector(`.title`);
    title.innerHTML = device.server.name || device.server.addr;

    /** @type {HTMLElement} */
    const editButton = item.querySelector(`button.edit`);
    editButton.setAttribute("data-addr", device.server.addr);

    /** @type {HTMLButtonElement} */
    const powerButton = item.querySelector(`button.power-button`);

    // @ts-expect-error
    powerButton.onclick = onClickPowerButton;

    // @ts-expect-error
    powerButton.querySelector(`.background`).style.backgroundColor =
        `rgb(${device.color.slice(0, 3).join(", ")})`;

    if (Math.max(...device.color)) {
        powerButton.setAttribute("data-state", "on");
    } else {
        powerButton.setAttribute("data-state", "off");
    }
}

/**
 * @param {Event & { currentTarget: HTMLButtonElement }} ev
 * @returns {Promise<void>}
 */
async function onClickPowerButton(ev) {
    const target = ev.currentTarget;

    if (target.getAttribute("data-state") === "processing") return;
    target.setAttribute("data-state", "processing");

    const defer = () => {
        if (device && Math.max(...device.color)) {
            target.setAttribute("data-state", "on");
        } else {
            target.setAttribute("data-state", "off");
        }
    };

    /** @type {string} */
    const addr = ev.currentTarget
        .closest(".device-list-item")
        .getAttribute("data-addr");

    // Search the local storage for this device
    let device = w.store.device(addr);

    // Set color
    /** @type {import("../types.d.ts").Color} */
    let newColor;
    if (!device.color || !device.color.find((c) => c > 0)) {
        newColor = [255, 255, 255, 255];
    } else {
        newColor = [0, 0, 0, 0];
    }

    // Request to api
    try {
        device = (await w.api.setDevicesColor(newColor, device))[0];
    } catch (err) {
        console.error(err);
        alert(err); // TODO: Error handling, notification?
        return defer();
    }

    // Update storage
    w.store.obj.update("devices", (storeDevices) => {
        for (let x = 0; x < storeDevices.length; x++) {
            if (storeDevices[x].server.addr === device.server.addr) {
                storeDevices[x] = device;
            }
        }
        return storeDevices;
    });

    return defer();
}
