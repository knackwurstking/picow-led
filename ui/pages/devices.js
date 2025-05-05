(() => {
    const { createDeviceItem } = require("../components/device-item");

    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-expect-error
    const w = window;

    /**
     * @param {import("../types.d.ts").Device} device
     */
    function currentColorForDevice(device) {
        return (
            w.store.currentColor(device.server.addr) ||
            (device.pins || []).map(() => 255)
        );
    }

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
                    const item = createDeviceItem(device, async () => {
                        let color = currentColorForDevice(device);
                        if (Math.max(...device.color) > 0) {
                            color = color.map(() => 0);
                        }

                        try {
                            await w.api.setDevicesColor(color, device);
                        } catch (err) {
                            console.error(err);
                            alert(err); // TODO: Error handling, notification?
                        }
                    });

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
