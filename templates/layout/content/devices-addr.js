//{{ define "script-devices-addr" }}
/** @type {PageWindow} */
// @ts-ignore
const w = window;

/**
 * @returns {string}
 */
function getDeviceAddress() {
    return decodeURIComponent(location.pathname.split("/").reverse()[0]);
}

/**
 * @param {UIStore} store
 * @returns {void}
 */
function setupAppBar(store) {
    const items = w.utils.setupAppBarItems(
        "back-button",
        "online-indicator",
        "title",
    );

    items["back-button"].onclick = (ev) => {
        ev.preventDefault();
        location.pathname = `{{ .ServerPathPrefix }}/`;
    };

    const addr = getDeviceAddress();
    /** @type {Device | undefined} */
    const device = (store.get("devices") || []).find((device) => {
        return device.server.addr === addr;
    });

    items["title"].innerText = device ? device.server.name : "";
}

window.addEventListener("load", async () => {
    /** @type {UIStore} */
    const store = new w.ui.Store("picow-led");

    setupAppBar(store);
});
//{{ end }}
