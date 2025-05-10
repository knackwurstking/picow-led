/**
 * @returns {import("../../types").Api}
 */
export function create() {
    /**
     * @returns {Promise<import("../../types").Device[]>}
     */
    async function devices() {
        const url = getURL("/api/devices");

        try {
            const resp = await fetch(url);

            try {
                const data = await handleResponse(resp, url);
                window.store.obj.set("devices", data);
                return data;
            } catch (err) {
                console.error(`Handle fetch response for ${url}:`, err);
            }
        } catch (err) {
            console.error(`fetch ${url}:`, err);
        }

        return window.store.obj.get("devices");
    }

    /**
     * @param {import("../../types").Color | undefined | null} color
     * @param {import("../../types").Device[]} devices
     * @returns {Promise<import("../../types").Device[]>}
     */
    async function setDevicesColor(color, ...devices) {
        // TODO: Do the same thing like in devices and colors

        if (!color) {
            color = [255, 255, 255, 255];
        }

        const url = getURL("/api/devices/color");
        const data = { devices, color };
        console.debug(`POST "${url}":`, data);

        const resp = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(data),
        });

        /** @type {import("../../types").Device[]} */
        devices = await handleResponse(resp, url);
        return updateStoreDevices(devices);
    }

    /**
     * @returns {Promise<import("../../types").Colors>}
     */
    async function colors() {
        const url = getURL("/api/colors");

        try {
            const resp = await fetch(url);

            try {
                const data = await handleResponse(resp, url);
                window.store.obj.set("colors", data);
                return data;
            } catch (err) {
                console.error(`Handle fetch response for ${url}:`, err);
            }
        } catch (err) {
            console.error(`Fetch ${url}:`, err);
        }

        return window.store.obj.get("colors");
    }

    /**
     * @param {number} index
     * @returns {Promise<import("../../types").Color | undefined>}
     */
    async function color(index) {
        const url = getURL(`/api/colors/${index}`);

        try {
            const resp = await fetch(url);

            try {
                /** @type {import("../../types").Color} */
                const color = await handleResponse(resp, url);

                window.store.obj.update("colors", (colors) => {
                    return colors.map((c, i) => (i === index ? color : c));
                });

                return color;
            } catch (err) {
                console.error(`Handle fetch response for ${url}:`, err);
            }
        } catch (err) {
            console.error(`Fetch ${url}:`, err);
        }

        return window.store.obj.get("colors")[index]; // Could be undefined
    }

    /**
     * @param {number} index
     * @param {import("../../types").Color} color
     * @returns {Promise<void>}
     */
    async function setColor(index, color) {
        // TODO: Do the same thing like in devices and colors
        const url = getURL(`/api/colors/${index}`);

        const resp = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(color),
        });

        return handleResponse(resp, url);
    }

    return {
        devices,
        setDevicesColor,
        colors,
        color,
        setColor,
    };
}

/**
 * @param {string} path
 * @returns {string}
 */
function getURL(path) {
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
 * @param {import("../../types").Device[]} devices
 * @returns {import("../../types").Device[]}
 */
function updateStoreDevices(devices) {
    window.store.obj.update("devices", (storeDevices) => {
        /** @type {import("../../types").Device} */
        let storeDevice;

        for (let sI = 0; sI < storeDevices.length; sI++) {
            for (let i = 0; i < devices.length; i++) {
                storeDevice = storeDevices[sI];

                if (storeDevice.server.addr === devices[i].server.addr) {
                    storeDevices[sI] = devices[i];

                    // Store current color
                    if (Math.max(...storeDevice.color) > 0) {
                        window.store.obj.update(
                            "currentDeviceColors",
                            (data) => {
                                data[storeDevice.server.addr] =
                                    storeDevice.color;

                                return data;
                            },
                        );
                    }

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
                }
            }
        }

        return storeDevices;
    });

    return devices;
}
