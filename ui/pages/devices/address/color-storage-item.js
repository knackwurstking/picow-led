export const colorSeparator = ",";

/**
 * @param {string} v
 * @return {Color}
 */
export function splitDataColor(v) {
    return v.split(colorSeparator).map((c) => parseInt(c, 10));
}

/**
 * @param {number} index
 * @param {Color} color
 * @param {Object} [options]
 * @param {Device} [options.device]
 * @param {(color: Color) => void|Promise<void>} [options.onClick]
 * @param {(color: Color) => void|Promise<void>} [options.onChange]
 * @returns {HTMLElement}
 */
export function create(index, color, options) {
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

    return update(item, index, color, options);
}

/**
 * @param {HTMLElement} item
 * @param {number} index
 * @param {Color} color
 * @param {Object} options
 * @param {Device | null} [options.device]
 * @param {(color: Color) => void|Promise<void>} [options.onClick]
 * @param {(color: Color) => void|Promise<void>} [options.onChange]
 * @returns {HTMLElement}
 */
export function update(item, index, color, options) {
    if (color.length < 3) color = [...color, 0, 0, 0];
    color = color.slice(0, 3);
    item.style.color = `rgb(${color.join(", ")})`;
    item.setAttribute("data-color", `${color.join(colorSeparator)}`);

    if (options.onClick) {
        item.onclick = () => {
            options.onClick(color);
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

        window.api.setColor(index, color);
        if (options.device) {
            window.api.setDevicesColor(color, options.device);
            if (options.onChange) options.onChange(color);
        }

        update(item, index, color, options);
    };

    return item;
}
