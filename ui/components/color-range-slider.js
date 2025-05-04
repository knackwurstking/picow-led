/**
 * @param {string} title
 * @param {number} value
 * @param {(ev: Event & { currentTarget: HTMLInputElement }) => void|Promise<void>} onChange
 * @returns {import("../types.d.ts").Component}
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

    return {
        element: item,
        destroy: updateColorRangeSlider(item, title, value, onChange),
    };
}

/**
 * @param {HTMLElement} item
 * @param {string} title
 * @param {number} value
 * @param {(ev: Event & { currentTarget: HTMLInputElement }) => void|Promise<void>} onChange
 * @returns {import("ui").CleanUpFunction}
 */
export function updateColorRangeSlider(item, title, value, onChange) {
    /** @type {HTMLElement} */
    const titleElement = item.querySelector(`.title`);

    /** @type {HTMLElement} */
    const circle = item.querySelector(`.circle`);
    /** @type {HTMLElement} */
    const range = item.querySelector(`.range`);
    /** @type {HTMLInputElement} */
    const input = item.querySelector(`input[type="number"]`);

    /** @type {DOMRect} */
    let rR;
    /** @type {DOMRect} */
    let cR;
    /** @type {number} */
    let xM = 1;

    /** @type {boolean} */
    let noneSelectBackup = false;
    /** @type {string} */
    let touchActionBackup = "";

    /** @type {number | null} */
    let currentPointerID = null;

    const updateRects = () => {
        rR = range.getBoundingClientRect();
        cR = circle.getBoundingClientRect();
    };

    /**
     * @returns {Record<"trackWidth" | "maxRange" | "minRange" | "circleRadius", number>}
     */
    const calculations = () => {
        const trackWidth = rR.width - cR.width - xM * 2;
        const maxRange = (trackWidth - cR.width) / (trackWidth / 100); // color: 0, circles border width is 2px
        const minRange = 100 - (trackWidth - xM) / (trackWidth / 100); // color: 255

        const circleRadius = (cR.right - cR.left) / 2;

        return {
            trackWidth,
            maxRange,
            minRange,
            circleRadius,
        };
    };

    /** @param {PointerEvent} ev */
    const pointerMove = (ev) => {
        if (currentPointerID !== ev.pointerId) {
            return;
        }

        const c = calculations();

        let right =
            (c.trackWidth - (ev.clientX - rR.left - xM)) / (c.trackWidth / 100);

        if (right >= c.maxRange) {
            right = c.maxRange;
        } else if (right <= c.minRange) {
            right = c.minRange;
        }

        circle.style.right = `${right}%`;

        const min = (100 - c.maxRange) * (c.trackWidth / 100); // 0
        const max = (100 - c.minRange) * (c.trackWidth / 100); // 255
        const current = (100 - right) * (c.trackWidth / 100);
        input.value = `${Math.round(((current - min) / ((max - min) / 100)) * 2.55)}`;
    };

    const pointerEnd = (/** @type {PointerEvent} */ ev) => {
        if (currentPointerID !== ev.pointerId) {
            return;
        }

        currentPointerID = null;
        window.removeEventListener("pointermove", pointerMove);
        window.removeEventListener("pointerup", pointerEnd);

        if (!noneSelectBackup) {
            document.body.classList.remove("ui-none-select");
        }

        document.body.style.touchAction = touchActionBackup;
    };

    const pointerStart = (/** @type {PointerEvent} */ ev) => {
        if (currentPointerID === ev.pointerId) {
            return;
        }
        ev.preventDefault();

        currentPointerID = ev.pointerId;
        noneSelectBackup = document.body.classList.contains("ui-none-select");
        document.body.classList.add("ui-none-select");

        touchActionBackup = document.body.style.touchAction;
        document.body.style.touchAction = "none";

        updateRects();

        window.addEventListener("pointerup", pointerEnd);
        window.addEventListener("pointermove", pointerMove);
    };

    /** @param {PointerEvent & { currentTarget: HTMLInputElement } | null} ev */
    input.onchange = (ev) => {
        updateRects();
        const c = calculations();

        let value = parseInt(input.value || "0", 10);
        if (value < 0) {
            value = 0;
            input.value = value.toString();
        } else if (value > 255) {
            value = 255;
            input.value = value.toString();
        }

        circle.style.right = `${100 - (100 - (100 - c.maxRange) - c.minRange) / (255 / value) - cR.width / (c.trackWidth / 100)}%`;

        if (ev) onChange(ev);
    };

    circle.onpointerdown = pointerStart;

    titleElement.innerText = title;
    input.value = value.toString();

    // NOTE: This only works with a `setTimeout` for now it seems
    setTimeout(() => {
        // @ts-expect-error
        input.onchange();
    });

    return () => {
        window.removeEventListener("pointerup", pointerEnd);
        window.removeEventListener("pointermove", pointerMove);

        circle.onpointerdown = null;
        input.onchange = null;
    };
}
