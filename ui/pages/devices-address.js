(() => {
    const {
        colorSeparator,
        createColorStorageItem,
    } = require("../items/color-storage.js");

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
     * @returns {void}
     */
    function setupAppBar() {
        const device = getDevice();
        const items = w.utils.setupAppBarItems("online-indicator", "title");
        items["title"].innerText = device ? device.server.name : "";
    }

    async function setupColorStorage() {
        /** @type {HTMLElement} */
        const colorCacheContainer = document.querySelector(
            `.color-storage-container`,
        );
        colorCacheContainer.innerHTML = "";

        const colorCache = await w.api.colors();

        for (let x = 0; x < colorCache.length; x++) {
            const item = createColorStorageItem(
                x,
                colorCache[x],
                getDevice(),
                (color) => {
                    const colorString = color.join(colorSeparator);

                    Array.from(colorCacheContainer.children).forEach(
                        (child) => {
                            if (
                                child.getAttribute("data-color") === colorString
                            ) {
                                if (!child.classList.contains("active")) {
                                    child.classList.add("active");

                                    const color = child
                                        .getAttribute(`data-color`)
                                        .split(colorSeparator)
                                        .map((/** @type{string} */ c) =>
                                            parseInt(c, 10),
                                        );

                                    if (
                                        color.length < 4 &&
                                        color.filter((c) => c === color[0])
                                            .length === 3
                                    ) {
                                        color.push(color[0]); // NOTE: Just some workaround for auto the missing white value (4. Pin)
                                    }

                                    w.api.setDevicesColor(color, getDevice());
                                }
                            } else {
                                child.classList.remove("active");
                            }
                        },
                    );
                },
            );

            colorCacheContainer.appendChild(item);
        }
    }

    window.addEventListener("pageshow", async () => {
        setupAppBar();
        setupColorStorage();
    });
})();
