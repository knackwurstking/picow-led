(() => {
    const deviceItem = require("./device-item");

    /**
     * @param {Device} device
     * @returns {Color}
     */
    function currentColorForDevice(device) {
        return (
            window.store.currentDeviceColor(device.server.addr) ||
            (device.pins || []).map(() => 255)
        );
    }

    async function setupAppBar() {
        const items = window.utils.setupAppBarItems(
            "online-indicator",
            "title",
            "settings-button",
        );

        items["title"].innerText = "Devices";
    }

    async function setupDevicesList() {
        const devices = await window.api.devices();

        /** @type {HTMLElement} */
        const devicesList = document.querySelector("._content.devices > .list");
        devicesList.innerHTML = "";

        /** @param {Device} device */
        const onClick = async (device) => {
            /** @type {Color} */
            let color;
            if (Math.max(...device.color) > 0) {
                color = (device.pins || device.color).map(() => 0);
            } else {
                color = currentColorForDevice(device);
            }

            // TODO: The websocket message handler will handle the device item update
            await window.api.setDevicesColor(color, device);
        };

        devices.forEach((device) => {
            const item = deviceItem.create(device, () => {
                onClick(device);
            });

            devicesList.appendChild(item);
        });

        window.ws.events.addListener("device", (device) => {
            /** @type {HTMLElement} */
            let child;
            for (let x = 0; x < devices.length; x++) {
                if (devices[x].server.addr !== device.server.addr) {
                    continue;
                }

                // @ts-expect-error
                child = devicesList.children[x];

                deviceItem.update(child, device, () => {
                    onClick(device);
                });
            }
        });
    }

    window.addEventListener("pageshow", async () => {
        setupAppBar();
        setupDevicesList();
    });
})();
