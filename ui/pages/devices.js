//!{{ define "script-page-devices" }}
/** @type {import("../types.d.ts").PageWindow} */
// @ts-expect-error
const w = window;

/**
 * @returns {void}
 */
function setupAppBar() {
    const items = w.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"].innerText = "Devices";
}

window.addEventListener("load", async () => {
    /** @type {import("../types.d.ts").UIStore} */
    const store = new w.ui.Store("picow-led:");

    store.set("devices", [], true);
    store.listen(
        "devices",
        (devices) => {
            /** @type {HTMLElement} */
            const devicesList = document.querySelector(
                "._content.devices > .list",
            );
            devicesList.innerHTML = "";

            devices.forEach((device) => {
                const item = createDeviceListItem(device);
                devicesList.appendChild(item);
            });
        },
        true,
    );

    setupAppBar();

    w.api.devices().then((devices) => {
        // Fetch Devices from the api (if not offline)
        store.set("devices", devices);
    });
});

/**
 * @param {Event & { currentTarget: HTMLButtonElement }} ev
 * @returns {Promise<void>}
 */
async function onClickPowerButton(ev) {
    // Disable rapid fire clicks
    const target = ev.currentTarget;

    // Backup state
    const prevState = target.getAttribute("data-state");
    if (prevState === "processing") return;

    // Lock, prevent rapid fire clicking
    target.setAttribute("data-state", "processing");

    // Get the device list item belonging to this button
    const deviceListItem = ev.currentTarget.closest(".device-list-item");

    /** @type {string} */
    const addr = deviceListItem.getAttribute("data-addr");

    // Search the local storage for this device
    /** @type {import("../types.d.ts").Device | null} */
    let device = null;

    /** @type {import("../types.d.ts").UIStore} */
    const store = new w.ui.Store("picow-led:");
    const storeDevices = store.get("devices") || [];

    for (const storeDevice of storeDevices) {
        if (storeDevice.server.addr === addr) {
            device = storeDevice;
            break;
        }
    }

    if (device === null) {
        throw new Error(`device for address ${device.server.addr} not found`);
    }

    // Set color
    /** @type {import("../types.d.ts").Color} */ let newColor;
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
        target.setAttribute("data-state", prevState);
        return;
    }

    // Update storage
    store.update("devices", (storeDevices) => {
        for (let x = 0; x < storeDevices.length; x++) {
            if (storeDevices[x].server.addr === device.server.addr) {
                storeDevices[x] = device;
            }
        }
        return storeDevices;
    });

    // Update .device-list-item
    /** @type {HTMLElement | null} */
    const item = document.querySelector(
        `.device-list-item[data-addr="${device.server.addr}"]`,
    );
    if (!item) {
        throw new Error(`device-list-item for ${device.server.addr} not found`);
    }
    updateDeviceListItem(item, device);

    // Set power button state
    if (Math.max(...device.color)) {
        target.setAttribute("data-state", "on");
    } else {
        target.setAttribute("data-state", "off");
    }
}

/**
 * @param {import("../types.d.ts").Device} device
 * @returns {HTMLElement}
 */
function createDeviceListItem(device) {
    /** @type {HTMLTemplateElement} */
    const template = document.querySelector(
        `template[name="device-list-item"]`,
    );

    /** @type {HTMLElement} */
    const item = template.content
        .cloneNode(true)
        // @ts-expect-error
        .querySelector(".device-list-item");

    updateDeviceListItem(item, device);

    return item;
}

/**
 * @param {HTMLElement} item
 * @param {import("../types.d.ts").Device} device
 * @returns {void}
 */
function updateDeviceListItem(item, device) {
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
}
//!{{ end }}
