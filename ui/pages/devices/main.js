(() => {
    const { createDeviceItem } = require("./device-item");

    /**
     * @param {import("../../types").Device} device
     * @returns {import("../../types").Color}
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
        const items = window.utils.setupAppBarItems(
            "online-indicator",
            "title",
            "settings-button",
        );

        items["title"].innerText = "Devices";
    }

    window.addEventListener("pageshow", async () => {
        setupAppBar();

        const devices = await window.api.devices();

        /** @type {HTMLElement} */
        const devicesList = document.querySelector("._content.devices > .list");
        devicesList.innerHTML = "";

        /** @param {import("../../types").Device} device */
        const createItem = (device) => {
            const onClick = async () => {
                /** @type {import("../../types").Color} */
                let color;
                if (Math.max(...device.color) > 0) {
                    color = color.map(() => 0);
                } else {
                    color = currentColorForDevice(device);
                }

                const devices = await window.api.setDevicesColor(color, device);
                devices.forEach(createItem);
            };

            const item = createDeviceItem(device, onClick);
            devicesList.appendChild(item);
        };

        devices.forEach(createItem);
    });
})();
