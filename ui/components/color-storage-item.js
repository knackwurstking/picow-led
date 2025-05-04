export const colorSeparator = ",";

/**
 * @param {string} v
 * @return {import("../types.d.ts").Color}
 */
export function splitDataColor(v) {
    return v.split(colorSeparator).map((c) => parseInt(c, 10));
}

/**
 * @param {number} index
 * @param {import("../types").Color} color
 * @param {import("../types").Device} [device]
 * @param {(color: import("../types").Color) => void|Promise<void>} [onClick]
 * @returns {import("../types.d.ts").Component}
 */
export function createColorStorageItem(index, color, device, onClick) {
    /** @type {HTMLTemplateElement} */
    const t = document.querySelector(`template[name="color-storage-item"]`);
    if (!t) {
        throw new Error(
            `Nope, template with name "color-range-slider" not found`,
        );
    }

    /** @type {HTMLElement} */
    // @ts-expect-error
    const item = t.content.cloneNode(true).querySelector(`*`);
    return {
        element: item,
        destroy: updateColorStorageItem(item, index, color, device, onClick),
    };
}

/**
 * @param {HTMLElement} item
 * @param {number} index
 * @param {import("../types").Color} color
 * @param {import("../types").Device | null} [device]
 * @param {(color: import("../types").Color) => void|Promise<void>} [onClick]
 * @returns {null}
 */
export function updateColorStorageItem(item, index, color, device, onClick) {
    if (color.length < 3) color = [...color, 0, 0, 0];
    color = color.slice(0, 3);
    item.style.color = `rgb(${color.join(", ")})`;
    item.setAttribute("data-color", `${color.join(colorSeparator)}`);

    if (onClick) {
        item.onclick = () => {
            onClick(color);
        };
    } else item.onclick = null;

    item.querySelector(`input`).onchange = (
        /** @type {Event & { currentTarget: HTMLInputElement }} */ ev,
    ) => {
        const value = (ev.currentTarget.value || "#FFFFFF").slice(1);
        const color = [];
        for (let x = 0; x < value.length; x += 2) {
            color.push(parseInt(value.slice(x, x + 2), 16));
        }

        /** @type {import("../types").PageWindow} */
        // @ts-ignore
        const w = window;
        w.api.setColor(index, color);
        if (device) w.api.setDevicesColor(color, device);

        updateColorStorageItem(item, index, color, device, onClick);
    };

    return null;
}
