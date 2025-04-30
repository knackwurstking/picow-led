//{{ define "script-window-utils" }}
(() => {
    /** @type {Api} */
    // @ts-ignore
    const api = window.api;

    /** @type {UI} */
    // @ts-ignore
    const ui = window.ui;

    /**
     * @param {Event & { currentTarget: HTMLButtonElement }} ev
     * @returns {Promise<void>}
     */
    async function onClickPowerButton(ev) {
        // Disable rapid fire clicks
        const target = ev.currentTarget;

        // Backup state
        const prevState = target.getAttribute("data-state");
        if (prevState === "processing") return;

        // Lock, prevent rapid fire clicking
        target.setAttribute("data-state", "processing");

        // Get the device list item belonging to this button
        const deviceListItem = ev.currentTarget.closest(".device-list-item");

        /** @type {string} */
        const addr = deviceListItem.getAttribute("data-addr");

        // Search the local storage for this device
        /** @type {Device | null} */
        let device = null;

        /** @type {UIStore} */
        const store = new ui.Store("picow-led");
        const storeDevices = store.get("devices") || [];

        for (const storeDevice of storeDevices) {
            if (storeDevice.server.addr === addr) {
                device = storeDevice;
                break;
            }
        }

        if (device === null) {
            throw new Error(
                `device for address ${device.server.addr} not found`,
            );
        }

        // Set color
        /** @type {Color} */
        let newColor;
        if (!device.color || !device.color.find((c) => c > 0)) {
            newColor = [255, 255, 255, 255];
        } else {
            newColor = [0, 0, 0, 0];
        }

        // Request to api
        try {
            device = (await api.setDevicesColor(newColor, device))[0];
        } catch (err) {
            console.error(err);
            alert(err); // TODO: Error handling, notification?
            target.setAttribute("data-state", prevState);
            return;
        }

        // Update storage
        store.update("devices", (storeDevices) => {
            for (let x = 0; x < storeDevices.length; x++) {
                if (storeDevices[x].server.addr === device.server.addr) {
                    storeDevices[x] = device;
                }
            }
            return storeDevices;
        });

        // Update .device-list-item
        /** @type {HTMLElement | null} */
        const item = document.querySelector(
            `.device-list-item[data-addr="${device.server.addr}"]`,
        );
        if (!item) {
            throw new Error(
                `device-list-item for ${device.server.addr} not found`,
            );
        }
        updateDeviceListItem(item, device);

        // Set power button state
        if (Math.max(...device.color)) {
            target.setAttribute("data-state", "on");
        } else {
            target.setAttribute("data-state", "off");
        }
    }

    /**
     * @param {HTMLElement} item
     * @param {Device} device
     * @returns {void}
     */
    function updateDeviceListItem(item, device) {
        item.setAttribute("data-addr", device.server.addr);

        /** @type {HTMLElement} */
        const title = item.querySelector(`.title`);
        title.innerHTML = device.server.name || device.server.addr;

        /** @type {HTMLElement} */
        const editButton = item.querySelector(`button.edit`);
        editButton.setAttribute("data-addr", device.server.addr);

        /** @type {HTMLButtonElement} */
        const powerButton = item.querySelector(`button.power-button`);

        // @ts-ignore
        powerButton.onclick = utils.onClickPowerButton;

        // @ts-ignore
        powerButton.querySelector(`.background`).style.backgroundColor =
            `rgb(${device.color.slice(0, 3).join(", ")})`;
    }

    /**
     * @param {HTMLElement} item
     * @param {string} name
     * @param {Color} color
     * @returns {void}
     */
    function updateColorCacheItem(item, name, color) {
        if (color.length < 3) color = [...color, 0, 0, 0];
        color = color.slice(0, 3);
        item.title = name;
        item.style.color = `rgb(${color.join(", ")})`;
    }

    /**
     * @param {AppBarItemName[]} itemNames
     * @returns {AppBarItems}
     */
    function setupAppBarItems(...itemNames) {
        /** @type {AppBarItems} */
        const enabledItems = {};

        /** @type {NodeListOf<HTMLElement>} */
        const items = document.querySelectorAll(`.ui-app-bar [data-name]`);
        let match = false;
        for (const item of items) {
            /** @type {AppBarItemName} */
            // @ts-ignore
            const dataName = item.getAttribute("data-name") || "";

            match = false;
            for (const name of itemNames) {
                if (name === dataName) {
                    match = true;
                }
            }

            if (match) {
                // Enable
                item.style.display = "inline-flex";
                enabledItems[dataName] = item;
            } else {
                // Disable
                item.style.display = "none";
            }
        }

        return enabledItems;
    }

    /**
     * @param {boolean} state
     * @returns {void}
     */
    function setOnlineIndicatorState(state) {
        const el = document.querySelector(`.online-indicator`);

        if (state) {
            el.setAttribute(`data-state`, "online");
        } else {
            el.setAttribute(`data-state`, "offline");
        }
    }

    /**
     * @returns {void}
     */
    function registerServiceWorker() {
        // Check if the browser supports service workers, otherwise abort.
        if (!("serviceWorker" in navigator)) {
            console.warn("Browser doesn't support service workers");
            return;
        }

        window.addEventListener("load", function () {
            navigator.serviceWorker
                .register("{{ .ServerPathPrefix }}/service-worker.js")
                .then(function (reg) {
                    console.info("Service worker registered", reg);
                })
                .catch(function (err) {
                    console.error("Service worker registration failed:", err);
                });
        });
    }

    /** @type {Utils} */
    const utils = {
        onClickPowerButton,
        updateDeviceListItem,
        updateColorCacheItem,
        setupAppBarItems,
        setOnlineIndicatorState,
        registerServiceWorker,
    };

    // @ts-ignore
    window.utils = utils;
})();
//{{ end }}
