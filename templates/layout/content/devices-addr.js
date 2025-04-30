//{{ define "script-devices-addr" }}
window.addEventListener("load", async () => {
    /** @type {Utils} */
    // @ts-ignore
    const utils = window.utils;

    /** @type {UI} */
    // @ts-ignore
    const ui = window.ui;

    // Setup AppBar Items

    const items = utils.setupAppBarItems(
        "back-button",
        "online-indicator",
        "title",
    );

    items["back-button"].onclick = (ev) => {
        ev.preventDefault();
        location.pathname = `{{ .ServerPathPrefix }}/`;
    };

    const addr = decodeURIComponent(location.pathname.split("/").reverse()[0]);

    /** @type {UIStore} */
    const store = new ui.Store("picow-led");

    /** @type {Device | undefined} */
    const device = (store.get("devices") || []).find((device) => {
        return device.server.addr === addr;
    });

    items["title"].innerText = device ? device.server.name : "";
});
//{{ end }}
