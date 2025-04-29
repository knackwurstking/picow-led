//{{ define "script-devices" }}
window.addEventListener("load", async () => {
    /** @type {UI} */
    // @ts-ignore
    const ui = window.ui;

    /** @type {UIStore} */
    const store = new ui.Store("picow-led");
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
                utils.updateDeviceListItem(item, device);
            });
        },
        true,
    );

    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    /** @type {Api} */
    // @ts-ignore
    const api = window.api;

    // Setup AppBar Items

    const items = utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"].innerText = "Devices";

    // Fetch Devices from the api (if not offline)

    api.devices().then((devices) => {
        store.set("devices", devices);
    });
});
//{{ end }}
