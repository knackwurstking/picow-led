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

        /** @param {Device} device */
        const createItem = (device) => {
            const onClick = async () => {
                /** @type {Color} */
                let color;
                if (Math.max(...device.color) > 0) {
                    color = color.map(() => 0);
                } else {
                    color = currentColorForDevice(device);
                }

                const devices = await window.api.setDevicesColor(color, device);
                devices.forEach(createItem);
            };

            const item = deviceItem.create(device, onClick);
            devicesList.appendChild(item);
        };

        devices.forEach(createItem);
    });
})();
