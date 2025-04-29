//{{ define "script-devices" }}
window.addEventListener("load", async () => {
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

    // Setup Devices List

    /** @type {HTMLElement} */
    const devicesList = document.querySelector("._content.devices > .list");
    devicesList.innerHTML = "";

    /** @type {HTMLTemplateElement} */
    const template = document.querySelector(
        `template[name="device-list-item"]`,
    );
    (await api.devices()).forEach((device) => {
        /** @type {HTMLElement} */
        const item = template.content
            .cloneNode(true)
            // @ts-ignore
            .querySelector(".device-list-item");

        devicesList.appendChild(item);
        utils.updateDeviceListItem(device);
    });
});
//{{ end }}
