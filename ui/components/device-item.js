/** @type {import("../types").PageWindow} */
// @ts-ignore
const w = window;

/**
 * @param {import("../types").Device} device
 * @param {(item: HTMLElement) => void|Promise<void>} onClickPowerButton
 * @returns {HTMLElement}
 */
export function createDeviceItem(device, onClickPowerButton) {
    /** @type {HTMLTemplateElement} */
    const t = document.querySelector(`template[name="device-list-item"]`);
    if (!t) {
        throw new Error(
            `Nope, template with name "color-range-slider" not found`,
        );
    }

    /** @type {HTMLElement} */
    const item = t.content
        .cloneNode(true)
        // @ts-expect-error
        .querySelector(".device-list-item");

    return updateDeviceItem(item, device, onClickPowerButton);
}

/**
 * @param {HTMLElement} item
 * @param {import("../types").Device} device
 * @param {(item: HTMLElement) => void|Promise<void>} onClickPowerButton
 * @returns {HTMLElement}
 */
export function updateDeviceItem(item, device, onClickPowerButton) {
    item.setAttribute("data-addr", device.server.addr);

    /** @type {HTMLElement} */
    const title = item.querySelector(`.title`);
    title.innerHTML = device.server.name || device.server.addr;

    /** @type {HTMLElement} */
    const editButton = item.querySelector(`button.edit`);
    editButton.setAttribute("data-addr", device.server.addr);

    /** @type {HTMLButtonElement} */
    const powerButton = item.querySelector(`button.power-button`);

    powerButton.onclick = async (
        /** @type {MouseEvent & { currentTarget: HTMLButtonElement }} */ ev,
    ) => {
        if (ev.currentTarget.getAttribute("data-state") === "processing") {
            return;
        }
        ev.currentTarget.setAttribute("data-state", "processing");

        await onClickPowerButton(item);
    };

    // @ts-expect-error
    powerButton.querySelector(`.background`).style.backgroundColor =
        `rgb(${device.color.slice(0, 3).join(", ")})`;

    if (Math.max(...device.color)) {
        powerButton.setAttribute("data-state", "on");
    } else {
        powerButton.setAttribute("data-state", "off");
    }

    return item;
}
