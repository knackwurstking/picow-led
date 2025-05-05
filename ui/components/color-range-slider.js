/**
 * @param {string} title
 * @param {number} value
 * @param {(ev: Event & { currentTarget: HTMLInputElement }) => void|Promise<void>} onChange
 * @returns {HTMLElement}
 */
export function createColorRangeSlider(title, value, onChange) {
    /** @type {HTMLTemplateElement} */
    const t = document.querySelector(`template[name="color-range-slider"]`);
    if (!t) {
        throw new Error(
            `Nope, template with name "color-range-slider" not found`,
        );
    }

    // @ts-expect-error
    const item = t.content.cloneNode(true).querySelector("*");

    return updateColorRangeSlider(item, title, value, onChange);
}

/**
 * @param {HTMLElement} item
 * @param {string} title
 * @param {number} value
 * @param {(ev: Event & { currentTarget: HTMLInputElement }) => void|Promise<void>} onChange
 * @returns {HTMLElement}
 */
export function updateColorRangeSlider(item, title, value, onChange) {
    /** @type {HTMLElement} */
    const titleElement = item.querySelector(`.title`);
    titleElement.innerText = title;

    /** @type {HTMLInputElement} */
    const rangeInput = item.querySelector(`input[type="range"]`);
    rangeInput.value = value.toString();
    rangeInput.oninput = () => {
        numberInput.value = rangeInput.value;
    };
    rangeInput.onchange = (
        /** @type {Event & { currentTarget: HTMLInputElement }} */ ev,
    ) => {
        onChange(ev);
    };

    /** @type {HTMLInputElement} */
    const numberInput = item.querySelector(`input[type="number"]`);
    numberInput.value = value.toString();
    numberInput.onchange = (
        /** @type {Event & { currentTarget: HTMLInputElement }} */ ev,
    ) => {
        rangeInput.value = numberInput.value;
        onChange(ev);
    };

    return item;
}
