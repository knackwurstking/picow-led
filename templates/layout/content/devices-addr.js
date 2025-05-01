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

    for (const name in colorCache) {
        const item = createColorCacheItem(name, colorCache[name], (color) => {
            const colorString = color.join(",");

            Array.from(colorCacheContainer.children).forEach((child) => {
                if (child.getAttribute("data-color") === colorString) {
                    child.classList.add("active");
                } else {
                    child.classList.remove("active");
                }
            });
        });

        colorCacheContainer.appendChild(item);
    }
}

window.addEventListener("load", async () => {
    /** @type {UIStore} */
    const store = new w.ui.Store("picow-led");

    setupAppBar(store);
    setupColorCache();
});

/**
 * @returns {HTMLElement}
 */
function createColorCacheItem(name, color, onClick) {
    /** @type {HTMLTemplateElement} */
    const t = document.querySelector(`template[name="color-cache-item"]`);
    /** @type {HTMLElement} */
    const item = t.content.cloneNode(true).querySelector(`*`);
    updateColorCacheItem(item, name, color, onClick);
    return item;
}

/**
 * @param {HTMLElement} item
 * @param {string} name
 * @param {Color} color
 * @param {(color: Color) => void|Promise<void>} onClick
 * @returns {void}
 */
function updateColorCacheItem(item, name, color, onClick) {
    if (color.length < 3) color = [...color, 0, 0, 0];
    color = color.slice(0, 3);
    item.style.color = `rgb(${color.join(", ")})`;
    item.setAttribute("data-color", `${color.join(",")}`);

    item.title = name;

    if (onClick) {
        item.onclick = () => {
            onClick(color);
        };
    } else item.onclick = null;

    const input = item.querySelector(`input`);
    input.onchange = () => {
        const value = (input.value || "#FFFFFF").slice(1);
        const color = [];
        for (let x = 0; x < value.length; x += 2) {
            color.push(parseInt(value.slice(x, x + 2), 16));
        }

        updateColorCacheItem(item, name, color, onClick);
    };
}
//{{ end }}
