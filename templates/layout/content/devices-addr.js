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

async function setupColorCache() {
    const colorCacheContainer = document.querySelector(
        `.color-cache-container`,
    );
    colorCacheContainer.innerHTML = "";

    const colorCache = await w.api.color();

    /** @type {HTMLTemplateElement} */
    const template = document.querySelector(
        `template[name="color-cache-item"]`,
    );

    for (const name in colorCache) {
        /** @type {HTMLElement} */
        const item = template.content
            .cloneNode(true)
            // @ts-ignore
            .querySelector(`.color-cache-item`);

        colorCacheContainer.appendChild(item);
        w.utils.updateColorCacheItem(item, name, colorCache[name], (color) => {
            // TODO: Update color picker here
        });
    }
}

window.addEventListener("load", async () => {
    /** @type {UIStore} */
    const store = new w.ui.Store("picow-led");

    setupAppBar(store);
    setupColorCache();
});
//{{ end }}
