(() => {
    const {
        colorSeparator,
        splitDataColor,
        createColorStorageItem,
    } = require("./color-storage-item.js");

    const { createColorRangeSlider } = require("./color-range-slider.js");

    /**
     * @returns {string}
     */
    function pageDeviceAddress() {
        return decodeURIComponent(location.pathname.split("/").reverse()[0]);
    }

    /**
     * @returns {Device}
     */
    function pageDevice() {
        return window.store.device(pageDeviceAddress());
    }

    /**
     * @returns {Color}
     */
    function pageCurrentColor() {
        return (
            window.store.currentDeviceColor(pageDeviceAddress()) ||
            (pageDevice().pins || []).map(() => 255)
        );
    }

    /**
     * @returns {Color | null}
     */
    function pagePickedColor() {
        // Get color from active item
        const color = pageActiveColor().slice(0, 3);

        // Get range slider values
        color.push(...pageRangeSliderValues());

        return color;
    }

    /**
     * @returns {Color}
     */
    function pageActiveColor() {
        let color = [];
        const activeItem = document.querySelector(`.color-storage-item.active`);
        if (activeItem) {
            color.push(
                ...splitDataColor(activeItem.getAttribute("data-color")),
            );
        } else {
            color = [255, 255, 255];
        }
        return color;
    }

    /**
     * @returns {number[]}
     */
    function pageRangeSliderValues() {
        return Array.from(
            document.querySelectorAll(
                `.range-sliders .color-range-slider input[type="range"]`,
            ),
        ).map((/**@type {HTMLInputElement}*/ input) => {
            return parseInt(input.value || "0", 10);
        });
    }

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const device = pageDevice();
        const items = window.utils.setupAppBarItems(
            "online-indicator",
            "title",
        );
        items["title"].innerText = device ? device.server.name : "";
    }

    async function setupColorStorage() {
        /** @type {HTMLElement} */
        const colorStorageContainer = document.querySelector(
            `.color-storage-container`,
        );
        colorStorageContainer.innerHTML = "";

        const currentColor = pageCurrentColor();
        const currentColorString = currentColor
            .slice(0, 3)
            .join(colorSeparator);

        // Create color storage items
        const device = pageDevice();
        (await window.api.colors()).forEach((color, index) => {
            /**
             * @param {Color} color
             */
            const onClick = async (color) => {
                const colorString = color.join(colorSeparator);

                Array.from(colorStorageContainer.children).forEach((child) => {
                    if (child.getAttribute("data-color") === colorString) {
                        if (!child.classList.contains("active")) {
                            child.classList.add("active");

                            const color = splitDataColor(
                                child.getAttribute("data-color"),
                            );

                            window.api.setDevicesColor(
                                [...color, ...pageRangeSliderValues()],
                                device,
                            );
                        }
                    } else {
                        child.classList.remove("active");
                    }
                });
            };

            // TODO: Add on change callback => do a api.setColor(...)
            const item = createColorStorageItem(index, color, device, onClick);

            if (item.getAttribute("data-color") === currentColorString) {
                item.classList.add("active");
            }

            colorStorageContainer.appendChild(item);
        });
    }

    async function setupRangeSliders() {
        /** @type {HTMLElement} */
        const container = document.querySelector(".range-sliders");
        container.innerHTML = "";

        const device = pageDevice();

        if (device.pins.length > 3) {
            container.style.display = "block";
        } else {
            container.style.display = "none";
            return;
        }

        if (device.pins) {
            const currentColor = pageCurrentColor();
            let timeout = null;
            device.pins.slice(3).forEach((pin, index) => {
                index += 3;
                const slider = createColorRangeSlider(
                    `Pin: ${pin.toString()}`,
                    currentColor[index] || 0,
                    () => {
                        // NOTE: Update device color (api) with some timeout
                        //       (250ms?), i should use websockets for this later
                        if (timeout !== null) {
                            clearTimeout(timeout);
                            timeout = null;
                        }
                        timeout = setTimeout(() => {
                            timeout = null;
                            window.api.setDevicesColor(
                                pagePickedColor(),
                                device,
                            );
                        }, 250);
                    },
                );

                container.appendChild(slider);
            });
        }
    }

    window.addEventListener("pageshow", async () => {
        setupAppBar();
        setupColorStorage();
        setupRangeSliders();

        console.debug("device address:", pageDevice());
        console.table({
            activeColor: pageActiveColor(),
            currentColor: pageCurrentColor(),
            pickedColor: pagePickedColor(),
            rangeSliderValues: pageRangeSliderValues(),
        });
    });
})();
