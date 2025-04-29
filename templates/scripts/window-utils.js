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
    async function powerButtonClickHandler(ev) {
        // Disable rapid fire clicks
        const target = ev.currentTarget;
        const prevState = target.getAttribute("data-state");
        if (prevState === "processing") return;
        target.setAttribute("data-state", "processing");

        const deviceListItem = ev.currentTarget.closest(".device-list-item");
        /** @type {Device} */
        let device = JSON.parse(deviceListItem.getAttribute("data-json"));

        /** @type {MicroColor} */
        let color;
        if (!device.color || !device.color.find((c) => c > 0)) {
            color = [255, 255, 255, 255];
        } else {
            color = [0, 0, 0, 0];
        }

        try {
            device = (await api.setDevicesColor(color, device))[0];
        } catch (err) {
            console.error(err);
            alert(err); // TODO: Error handling, notification?
            target.setAttribute("data-state", prevState);
            return;
        }

        // @ts-ignore
        deviceListItem.querySelector(".title").innerText =
            device.server.name || device.server.addr;

        deviceListItem.setAttribute("data-json", JSON.stringify(device));
        if (Math.max(...device.color) > 0) {
            target.setAttribute("data-state", "on");
        } else {
            target.setAttribute("data-state", "off");
        }

        /** @type {HTMLElement} */
        const bg = deviceListItem.querySelector("div.background");
        bg.style.backgroundColor = `rgb(${device.color.slice(0, 3).join(", ")})`;
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
        for (const item of items) {
            const dataName = item.getAttribute("data-name") || "";

            for (const name of itemNames) {
                if (name === dataName) {
                    // Enable
                    item.style.display = "inline-flex";
                    enabledItems[name] = item;
                } else {
                    // Disable
                    item.style.display = "none";
                }
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
        powerButtonClickHandler,
        setupAppBarItems,
        setOnlineIndicatorState,
        registerServiceWorker,
    };

    // @ts-ignore
    window.utils = utils;
})();
//{{ end }}
