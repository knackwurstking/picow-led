(() => {
    const { createDeviceItem } = require("../items/device");

    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-expect-error
    const w = window;

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = w.utils.setupAppBarItems(
            "online-indicator",
            "title",
            "settings-button",
        );

        items["title"].innerText = "Devices";
    }

    window.addEventListener("pageshow", async () => {
        w.store.obj.listen(
            "devices",
            (devices) => {
                /** @type {HTMLElement} */
                const devicesList = document.querySelector(
                    "._content.devices > .list",
                );
                devicesList.innerHTML = "";

                devices.forEach((device) => {
                    const item = createDeviceItem(device);
                    devicesList.appendChild(item);
                });
            },
            true,
        );

        setupAppBar();

        w.api.devices().then((devices) => {
            // Fetch Devices from the api (if not offline)
            w.store.obj.set("devices", devices);
        });
    });
})();
