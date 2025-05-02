(() => {
    /**
     * @returns {import("./types.d.ts").UIStore}
     */
    function createStore() {
        /** @type {import("./types.d.ts").PageWindow} */
        // @ts-expect-error
        const w = window;

        /** @type {import("./types.d.ts").UIStore} */
        const store = new w.ui.Store("picow-led:");

        store.set("devices", [], true);

        return store;
    }

    /**
     * @returns {import("./types.d.ts").Api}
     */
    function createApi() {
        /**
         * @param {string} path
         * @returns {string}
         */
        function getUrl(path) {
            // @ts-ignore
            return process.env.SERVER_PATH_PREFIX + `${path}`;
        }

        /**
         * @param {Response} resp
         * @param {string} url
         * @returns {Promise<any>}
         */
        async function _handleResponse(resp, url) {
            const status = resp.status;

            if (!resp.ok) {
                throw new Error(`${status}: ${(await resp.text()) || "???"}`);
            }

            const respData = await resp.json();
            console.debug(`Got data from "${url}":`, respData);
            return respData;
        }

        /**
         * @returns {Promise<import("./types.d.ts").Device[]>}
         */
        async function devices() {
            const url = getUrl("/api/devices");

            const resp = await fetch(url);
            return await _handleResponse(resp, url);
        }

        /**
         * @param {import("./types.d.ts").Color | undefined | null} color
         * @param {import("./types.d.ts").Device[]} devices
         * @returns {Promise<import("./types.d.ts").Device[]>}
         */
        async function setDevicesColor(color, ...devices) {
            if (!color) {
                color = [255, 255, 255, 255];
            }

            const url = getUrl("/api/devices/color");
            const data = { devices, color };
            console.debug(`POST "${url}":`, data);

            const resp = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            });
            return _handleResponse(resp, url);
        }

        /**
         * @returns {Promise<import("./types.d.ts").ColorCache>}
         */
        async function colors() {
            const url = getUrl("/api/color");
            const resp = await fetch(url);
            return _handleResponse(resp, url);
        }

        /**
         * @param {number} index
         * @returns {Promise<import("./types.d.ts").Color>}
         */
        async function color(index) {
            const url = getUrl(`/api/color/${index}`);
            const resp = await fetch(url);
            return _handleResponse(resp, url);
        }

        /**
         * @param {number} index
         * @param {import("./types.d.ts").Color} color
         * @returns {Promise<void>}
         */
        async function setColor(index, color) {
            const url = getUrl(`/api/color/${index}`);

            const resp = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(color),
            });

            return _handleResponse(resp, url);
        }

        return {
            devices,
            setDevicesColor,
            colors,
            color,
            setColor,
        };
    }

    /**
     * @returns {import("./types.d.ts").WS}
     */
    function createWS() {
        /** @type {WebSocket | null} */
        let socket = null;
        /** @type {number | null} */
        let timeout = null;
        /** @type {number} */
        const timeoutDuration = 1000;

        const onClose = function () {
            if (timeout !== null) {
                clearTimeout(timeout);
                timeout = null;
            }

            // Reconnect here
            timeout = setTimeout(() => {
                connect();
            }, timeoutDuration);
        };

        const onOpen = () => {
            // TODO: Request devices
            //const devices = window.api.devices();
        };

        const onMessage = () => {
            // TODO: ...
        };

        function addr() {
            return ``; // TODO: ...
        }

        function isOpen() {
            if (!socket) return false;
            return socket.readyState === socket.OPEN;
        }

        function connect() {
            if (socket) close();

            const wsAddr = addr(); // origin + path
            console.debug(`Try to connect WebSocket to ${wsAddr}`);

            socket = new WebSocket(wsAddr);

            // Reconnect handler
            socket.addEventListener("close", onClose);
            socket.addEventListener("open", onOpen);
            socket.addEventListener("message", onMessage);
        }

        function close() {
            if (timeout) {
                clearTimeout(timeout);
                timeout = null;
            }

            if (socket) {
                socket.removeEventListener("close", onClose);
                if (isOpen()) socket.close();
                socket = null;
            }
        }

        /** @type {import("./types.d.ts").WS} */
        return {
            addr,
            isOpen,
            connect,
            close,
        };
    }

    /**
     * @returns {import("./types.d.ts").Utils}
     */
    function createUtils() {
        /**
         * @param {import("./types.d.ts").AppBarItemName[]} itemNames
         * @returns {import("./types.d.ts").AppBarItems}
         */
        function setupAppBarItems(...itemNames) {
            /** @type {import("./types.d.ts").AppBarItems} */
            const enabledItems = {};

            /** @type {NodeListOf<HTMLElement>} */
            const items = document.querySelectorAll(`.ui-app-bar [data-name]`);
            let match = false;
            for (const item of items) {
                /** @type {import("./types.d.ts").AppBarItemName} */
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
            if (!el) return;

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
                    .register(
                        // @ts-ignore
                        process.env.SERVER_PATH_PREFIX + "/service-worker.js",
                    )
                    .then(function (reg) {
                        console.info("Service worker registered", reg);
                    })
                    .catch(function (err) {
                        console.error(
                            "Service worker registration failed:",
                            err,
                        );
                    });
            });
        }

        return {
            setupAppBarItems,
            setOnlineIndicatorState,
            registerServiceWorker,
        };
    }

    // @ts-ignore
    window.store = createStore();
    // @ts-ignore
    window.api = createApi();
    // @ts-ignore
    window.ws = createWS();
    // @ts-ignore
    window.utils = createUtils();
})();
