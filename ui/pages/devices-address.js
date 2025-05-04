const {
    colorSeparator,
    createColorStorageItem,
} = require("../components/color-storage-item.js");

const {
    createColorRangeSlider,
} = require("../components/color-range-slider.js");

/** @type {import("../types.d.ts").PageWindow} */
// @ts-expect-error
const w = window;

/**
 * @returns {string}
 */
function getDeviceAddress() {
    return decodeURIComponent(location.pathname.split("/").reverse()[0]);
}

/**
 * @returns {import("../types.d.ts").Device}
 */
function getDevice() {
    return w.store.device(getDeviceAddress());
}

/**
 * @param {import("../types.d.ts").Device} device
 * @returns {import("../types.d.ts").Color | null}
 */
function getColor(device) {
    /** @type {import("../types.d.ts").Color} */
    const color = [];

    // Get color from active item, or use device color as fallback, or use [255,255,255] as default
    /** @type {HTMLElement | null} */
    const activeItem = document.querySelector(`.color-storage-item.active`);
    if (!activeItem) {
        color.push(
            // Ok, please don't ask questions here
            ...[...(device.color || [255, 255, 255]), 0, 0, 0].slice(0, 3),
        );
    }

    // Get range slider values
    document
        .querySelectorAll(".range-sliders .color-range-slider input")
        .forEach((/**@type {HTMLInputElement}*/ input) => {
            color.push(parseInt(input.value || "0", 10));
        });

    return [];
}

/**
 * @returns {void}
 */
function setupAppBar() {
    const device = getDevice();
    const items = w.utils.setupAppBarItems("online-indicator", "title");
    items["title"].innerText = device ? device.server.name : "";
}

async function setupColorStorage() {
    /** @type {HTMLElement} */
    const colorStorageContainer = document.querySelector(
        `.color-storage-container`,
    );
    colorStorageContainer.innerHTML = "";

    const colorCache = await w.api.colors();

    // Create color storage items
    for (let x = 0; x < colorCache.length; x++) {
        const item = createColorStorageItem(
            x,
            colorCache[x],
            getDevice(),
            (color) => {
                const colorString = color.join(colorSeparator);

                Array.from(colorStorageContainer.children).forEach((child) => {
                    if (child.getAttribute("data-color") === colorString) {
                        if (!child.classList.contains("active")) {
                            child.classList.add("active");

                            const color = child
                                .getAttribute(`data-color`)
                                .split(colorSeparator)
                                .map((/** @type{string} */ c) =>
                                    parseInt(c, 10),
                                );

                            // TODO: This fix is no longer needed, instead update the range sliders?
                            if (
                                color.length < 4 &&
                                color.filter((c) => c === color[0]).length === 3
                            ) {
                                color.push(color[0]); // NOTE: Just some workaround for auto fixing the missing 4. value (white)
                            }

                            w.api.setDevicesColor(color, getDevice());
                        }
                    } else {
                        child.classList.remove("active");
                    }
                });
            },
        );

        colorStorageContainer.appendChild(item.element);
    }
}

async function setupRangeSliders() {
    /** @type {HTMLElement} */
    const container = document.querySelector(".range-sliders");
    container.innerHTML = "";

    const device = getDevice();

    if (device.pins.length > 3) {
        container.style.display = "block";
    } else {
        container.style.display = "none";
        return;
    }

    if (device.pins) {
        device.pins.slice(3).forEach((pin, index) => {
            index += 3;
            const slider = createColorRangeSlider(
                pin.toString(),
                device.color[index] || 0,
                () => {
                    // NOTE: Update device color (api) with some timeout
                    //       (250ms?), i should use websockets for this later
                    setTimeout(() => {
                        w.api.setDevicesColor(getColor(device), device);
                    }, 250);
                },
            );

            container.appendChild(slider.element);
        });
    }
}

window.addEventListener("pageshow", async () => {
    setupAppBar();
    setupColorStorage();
    setupRangeSliders();
});
