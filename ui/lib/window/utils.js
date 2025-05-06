/**
 * @returns {import("../../types").Utils}
 */
export function create() {
    /**
     * @param {import("../../types").AppBarItemName[]} itemNames
     * @returns {import("../../types").AppBarItems}
     */
    function setupAppBarItems(...itemNames) {
        /** @type {import("../../types").AppBarItems} */
        const enabledItems = {};

        /** @type {NodeListOf<HTMLElement>} */
        const items = document.querySelectorAll(`.ui-app-bar [data-name]`);
        let match = false;
        for (const item of items) {
            /** @type {import("../../types").AppBarItemName} */
            // @ts-expect-error
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

        window.addEventListener("pageshow", function () {
            navigator.serviceWorker
                .register(process.env.SERVER_PATH_PREFIX + "/service-worker.js")
                .then(function (reg) {
                    console.info("Service worker registered", reg);
                })
                .catch(function (err) {
                    console.error("Service worker registration failed:", err);
                });
        });
    }

    return {
        setupAppBarItems,
        setOnlineIndicatorState,
        registerServiceWorker,
    };
}
