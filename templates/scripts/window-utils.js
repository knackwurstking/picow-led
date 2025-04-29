//{{ define "script-window-utils" }}
(() => {
    /** @type {Api} */
    // @ts-ignore
    const api = window.api;

    /**
     * TODO: Out of Date
     *
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
        const addr = JSON.parse(deviceListItem.getAttribute("data-addr"));

        // TODO: Get the color for this device from the storage somehow
        /** @type {Device} */
        let device = {
            // @ts-ignore
            server: {
                addr: addr,
            },
        };

        /** @type {MicroColor} */
        let newColor;
        if (!device.color || !device.color.find((c) => c > 0)) {
            newColor = [255, 255, 255, 255];
        } else {
            newColor = [0, 0, 0, 0];
        }

        try {
            device = (await api.setDevicesColor(newColor, device))[0];
        } catch (err) {
            console.error(err);
            alert(err); // TODO: Error handling, notification?
            target.setAttribute("data-state", prevState);
            return;
        }

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

        if (Math.max(...device.color)) {
            target.setAttribute("data-state", "on");
        } else {
            target.setAttribute("data-state", "off");
        }

        // TODO: Update storage, but for now i don't have any storage for this
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
                .register("{{ .ServerPathPrefix }}/js/service-worker.js", {
                    scope: "{{ .ServerPathPrefix }}",
                })
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
        setupAppBarItems,
        setOnlineIndicatorState,
        registerServiceWorker,
    };

    // @ts-ignore
    window.utils = utils;
})();
//{{ end }}
