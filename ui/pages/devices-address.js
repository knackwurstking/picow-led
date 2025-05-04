const {
    colorSeparator,
    splitDataColor,
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
 * @returns {import("../types.d.ts").Color | null}
 */
function getColor() {
    // Get color from active item
    const color = getActiveColor().slice(0, 3);

    // Get range slider values
    color.push(...getRangeSliderValues());

    return color;
}

/**
 * @returns {import("../types.d.ts").Color}
 */
function getActiveColor() {
    let color = [];
    const activeItem = document.querySelector(`.color-storage-item.active`);
    if (activeItem) {
        color.push(...splitDataColor(activeItem.getAttribute("data-color")));
    } else {
        color = [255, 255, 255];
    }
    return color;
}

/**
 * @returns {number[]}
 */
function getRangeSliderValues() {
    return Array.from(
        document.querySelectorAll(".range-sliders .color-range-slider input"),
    ).map((/**@type {HTMLInputElement}*/ input) => {
        return parseInt(input.value || "0", 10);
    });
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

                            const color = splitDataColor(
                                child.getAttribute("data-color"),
                            );

                            w.api.setDevicesColor(
                                [...color, ...getRangeSliderValues()],
                                getDevice(),
                            );
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
        let timeout = null;
        device.pins.slice(3).forEach((pin, index) => {
            index += 3;
            const slider = createColorRangeSlider(
                `Pin: ${pin.toString()}`,
                device.color[index] || 0,
                () => {
                    // NOTE: Update device color (api) with some timeout
                    //       (250ms?), i should use websockets for this later
                    console.debug(`range slider input change event`);
                    if (timeout !== null) {
                        clearTimeout(timeout);
                        timeout = null;
                    }
                    timeout = setTimeout(() => {
                        timeout = null;
                        w.api.setDevicesColor(getColor(), device);
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
