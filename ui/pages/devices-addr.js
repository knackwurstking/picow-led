//!{{ define "script-page-devices-addr" }}
(() => {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-expect-error
    const w = window;

    const dataColorSeparator = ",";

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
        const addr = getDeviceAddress();

        return (w.store.get("devices") || []).find((device) => {
            return device.server.addr === addr;
        });
    }

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = w.utils.setupAppBarItems(
            "back-button",
            "online-indicator",
            "title",
        );

        items["back-button"].onclick = (ev) => {
            ev.preventDefault();
            location.pathname = `{{ .ServerPathPrefix }}/`;
        };

        const device = getDevice();
        items["title"].innerText = device ? device.server.name : "";
    }

    async function setupColorStorage() {
        /** @type {HTMLElement} */
        const colorCacheContainer = document.querySelector(
            `.color-storage-container`,
        );
        colorCacheContainer.innerHTML = "";

        const colorCache = await w.api.colors();

        for (const name in colorCache) {
            const item = createColorCacheItem(
                name,
                colorCache[name],
                (color) => {
                    const colorString = color.join(dataColorSeparator);

                    Array.from(colorCacheContainer.children).forEach(
                        (child) => {
                            if (
                                child.getAttribute("data-color") === colorString
                            ) {
                                if (!child.classList.contains("active")) {
                                    child.classList.add("active");

                                    const color = child
                                        .getAttribute(`data-color`)
                                        .split(dataColorSeparator)
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

    window.addEventListener("load", async () => {
        setupAppBar();
        setupColorStorage();
    });

    /**
     * @param {string} name
     * @param {import("../types.d.ts").Color} color
     * @param {(color: import("../types.d.ts").Color) => void|Promise<void>} onClick
     * @returns {HTMLElement}
     */
    function createColorCacheItem(name, color, onClick) {
        /** @type {HTMLTemplateElement} */
        const t = document.querySelector(`template[name="color-storage-item"]`);

        /** @type {HTMLElement} */
        // @ts-expect-error
        const item = t.content.cloneNode(true).querySelector(`*`);
        updateColorCacheItem(item, name, color, onClick);
        return item;
    }

    /**
     * @param {HTMLElement} item
     * @param {string} name
     * @param {import("../types.d.ts").Color} color
     * @param {(color: import("../types.d.ts").Color) => void|Promise<void>} onClick
     * @returns {void}
     */
    function updateColorCacheItem(item, name, color, onClick) {
        if (color.length < 3) color = [...color, 0, 0, 0];
        color = color.slice(0, 3);
        item.style.color = `rgb(${color.join(", ")})`;
        item.setAttribute("data-color", `${color.join(dataColorSeparator)}`);

        item.title = name;

        if (onClick) {
            item.onclick = () => {
                onClick(color);
            };
        } else item.onclick = null;

        const input = item.querySelector(`input`);
        input.onchange = () => {
            const value = (input.value || "#FFFFFF").slice(1);
            const color = [];
            for (let x = 0; x < value.length; x += 2) {
                color.push(parseInt(value.slice(x, x + 2), 16));
            }

            w.api.setColor(name, color);
            w.api.setDevicesColor(color, getDevice());

            updateColorCacheItem(item, name, color, onClick);
        };
    }
})();
//!{{ end }}
