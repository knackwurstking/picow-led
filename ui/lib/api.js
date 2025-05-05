/**
 * @returns {Promise<import("../types").Device[]>}
 */
export async function devices() {
    const url = getUrl("/api/devices");

    const resp = await fetch(url);
    return await handleResponse(resp, url);
}

/**
 * @param {import("../types").Color | undefined | null} color
 * @param {import("../types").Device[]} devices
 * @returns {Promise<import("../types").Device[]>}
 */
export async function setDevicesColor(color, ...devices) {
    if (!color) {
        color = [255, 255, 255, 255];
    }

    const url = getUrl("/api/devices/color");
    const data = { devices, color };
    console.debug(`POST "${url}":`, data);

    const resp = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    /** @type {import("../types").Device[]} */
    devices = await handleResponse(resp, url);
    return updateStoreDevices(devices);
}

/**
 * @returns {Promise<import("../types").ColorCache>}
 */
export async function colors() {
    const url = getUrl("/api/color");
    const resp = await fetch(url);
    return handleResponse(resp, url);
}

/**
 * @param {number} index
 * @returns {Promise<import("../types").Color>}
 */
export async function color(index) {
    const url = getUrl(`/api/color/${index}`);
    const resp = await fetch(url);
    return handleResponse(resp, url);
}

/**
 * @param {number} index
 * @param {import("../types").Color} color
 * @returns {Promise<void>}
 */
export async function setColor(index, color) {
    const url = getUrl(`/api/color/${index}`);

    const resp = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(color),
    });

    return handleResponse(resp, url);
}

/**
 * @param {string} path
 * @returns {string}
 */
function getUrl(path) {
    return process.env.SERVER_PATH_PREFIX + `${path}`;
}

/**
 * @param {Response} resp
 * @param {string} url
 * @returns {Promise<any>}
 */
async function handleResponse(resp, url) {
    const status = resp.status;

    if (!resp.ok) {
        throw new Error(`${status}: ${(await resp.text()) || "???"}`);
    }

    const respData = await resp.json();
    console.debug(`Got data from "${url}":`, respData);
    return respData;
}

/**
 * @param {import("../types").Device[]} devices
 * @returns {import("../types").Device[]}
 */
function updateStoreDevices(devices) {
    window.store.obj.update("devices", (storeDevices) => {
        /** @type {import("../types").Device} */
        let storeDevice;

        for (let sI = 0; sI < storeDevices.length; sI++) {
            for (let i = 0; i < devices.length; i++) {
                storeDevice = storeDevices[sI];

                if (storeDevice.server.addr === devices[i].server.addr) {
                    storeDevices[sI] = devices[i];

                    // Log device error
                    if (storeDevice.error) {
                        console.error(
                            `Device ${
                                storeDevice.server.name ||
                                storeDevice.server.addr
                            } is ${
                                storeDevice.online ? "online" : "offline"
                            } with error:`,
                            storeDevice.error,
                        );
                    }

                    // Store current color
                    if (Math.max(...storeDevice.color) > 0) {
                        window.store.obj.update("color", (data) => {
                            data.current[storeDevice.server.addr] =
                                storeDevice.color;
                            return data;
                        });
                    }
                }
            }
        }

        return storeDevices;
    });

    return devices;
}
