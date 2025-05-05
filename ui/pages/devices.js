(() => {
    const utils = require("../lib/utils");
    const api = require("../lib/api");
    const { createDeviceItem } = require("../lib/components/device-item");

    /**
     * @param {import("../types").Device} device
     */
    function currentColorForDevice(device) {
        return (
            window.store.currentColor(device.server.addr) ||
            (device.pins || []).map(() => 255)
        );
    }

    /**
     * @returns {void}
     */
    function setupAppBar() {
        const items = utils.setupAppBarItems(
            "online-indicator",
            "title",
            "settings-button",
        );

        items["title"].innerText = "Devices";
    }

    window.addEventListener("pageshow", async () => {
        window.store.obj.listen(
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
                            await api.setDevicesColor(color, device);
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

        api.devices().then((devices) => {
            // Fetch Devices from the api (if not offline)
            window.store.obj.set("devices", devices);
        });
    });
})();
