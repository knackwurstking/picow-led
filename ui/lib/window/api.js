/**
 * @returns {Api}
 */
export function create() {
    return {
        /**
         * @returns {Promise<Device[]>}
         */
        async devices() {
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
        },

        /**
         * @param {Color | undefined | null} color
         * @param {Device[]} devices
         * @returns {Promise<Device[]>}
         */
        async setDevicesColor(color, ...devices) {
            if (!color) {
                color = [255, 255, 255, 255];
            }

            const url = getURL("/api/devices/color");
            const data = { devices, color };
            console.debug(`POST "${url}":`, data);

            try {
                const resp = await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(data),
                });

                try {
                    /** @type {Device[]} */
                    devices = await handleResponse(resp, url);
                } catch (err) {
                    console.error(`Handle fetch response from ${url}:`, err);
                }
            } catch (err) {
                console.error(`Fetch ${url}:`, err);
            }

            return updateDevicesStore(devices);
        },

        /**
         * @returns {Promise<Colors>}
         */
        async colors() {
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
        },

        /**
         * @param {number} index
         * @returns {Promise<Color>}
         */
        async color(index) {
            const url = getURL(`/api/colors/${index}`);

            try {
                const resp = await fetch(url);

                try {
                    /** @type {Color} */
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

            const color = window.store.obj.get("colors")[index];
            if (!color) {
                return await this.colors()[index];
            }
            return color;
        },

        /**
         * @param {number} index
         * @param {Color} color
         * @returns {Promise<void>}
         */
        async setColor(index, color) {
            const url = getURL(`/api/colors/${index}`);

            try {
                const resp = await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(color),
                });

                try {
                    const data = await handleResponse(resp, url);

                    window.store.obj.update("colors", (colors) => {
                        return colors.map((c, i) => (i === index ? color : c));
                    });

                    return data;
                } catch (err) {
                    console.error(`Handle fetch response for ${url}:`, err);
                }
            } catch (err) {
                console.error(`Fetch ${url}:`, err);
            }
        },
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
 * @param {Device[]} devices
 * @returns {Device[]}
 */
function updateDevicesStore(devices) {
    window.store.obj.update("devices", (storeDevices) => {
        /** @type {Device} */
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
