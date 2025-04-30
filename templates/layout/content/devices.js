//{{ define "script-devices" }}
/** @type {PageWindow} */
// @ts-ignore
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

window.addEventListener("load", async () => {
    /** @type {UIStore} */
    const store = new w.ui.Store("picow-led");

    store.set("devices", [], true);
    store.listen(
        "devices",
        (devices) => {
            /** @type {HTMLElement} */
            const devicesList = document.querySelector(
                "._content.devices > .list",
            );
            devicesList.innerHTML = "";

            /** @type {HTMLTemplateElement} */
            const template = document.querySelector(
                `template[name="device-list-item"]`,
            );

            devices.forEach((device) => {
                /** @type {HTMLElement} */
                const item = template.content
                    .cloneNode(true)
                    // @ts-ignore
                    .querySelector(".device-list-item");

                devicesList.appendChild(item);
                w.utils.updateDeviceListItem(item, device);
            });
        },
        true,
    );

    setupAppBar();

    w.api.devices().then((devices) => {
        // Fetch Devices from the api (if not offline)
        store.set("devices", devices);
    });
});
//{{ end }}
