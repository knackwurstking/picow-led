export function create(): Api {
    return {
        async devices(): Promise<Device[]> {
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

            return window.store.obj.get("devices") || [];
        },

        async setDevicesColor(
            color?: Color | null,
            ...devices: Device[]
        ): Promise<void> {
            const url = getURL("/api/devices/color?force=true");

            if (!color) {
                color = [255, 255, 255, 255];
            }

            const data = { devices, color };
            console.debug(`POST "${url}":`, data);

            try {
                await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(data),
                });
            } catch (err) {
                console.error(`Fetch ${url}:`, err);
            }
        },

        async colors(): Promise<Colors> {
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

            return window.store.obj.get("colors") || [];
        },

        async color(index: number): Promise<Color> {
            const url = getURL(`/api/colors/${index}`);

            try {
                const resp = await fetch(url);

                try {
                    const color: Color = await handleResponse(resp, url);

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

            const color = (window.store.obj.get("colors") || [])[index];
            if (!color) {
                return (await this.colors())[index];
            }
            return color;
        },

        async setColor(index: number, color: Color): Promise<void> {
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

function getURL(path: string): string {
    return process.env.SERVER_PATH_PREFIX + `${path}`;
}

async function handleResponse(resp: Response, url: string): Promise<any> {
    const status = resp.status;

    if (!resp.ok) {
        throw new Error(`${status}: ${(await resp.text()) || "???"}`);
    }

    const respData = await resp.json();
    console.debug(`Got data from "${url}":`, respData);
    return respData;
}

///**
// * @param {Device[]} devices
// * @returns {Device[]}
// */
//function updateDevicesStore(devices) {
//    window.store.obj.update("devices", (storeDevices) => {
//        /** @type {Device} */
//        let storeDevice;
//
//        for (let sI = 0; sI < storeDevices.length; sI++) {
//            for (let i = 0; i < devices.length; i++) {
//                storeDevice = storeDevices[sI];
//
//                if (storeDevice.server.addr === devices[i].server.addr) {
//                    storeDevices[sI] = devices[i];
//
//                    // Store current color
//                    if (Math.max(...storeDevice.color) > 0) {
//                        window.store.obj.update(
//                            "currentDeviceColors",
//                            (data) => {
//                                data[storeDevice.server.addr] =
//                                    storeDevice.color;
//
//                                return data;
//                            },
//                        );
//                    }
//
//                    // Log device error
//                    if (storeDevice.error) {
//                        console.error(
//                            `Device ${
//                                storeDevice.server.name ||
//                                storeDevice.server.addr
//                            } is ${
//                                storeDevice.online ? "online" : "offline"
//                            } with error:`,
//                            storeDevice.error,
//                        );
//                    }
//                }
//            }
//        }
//
//        return storeDevices;
//    });
//
//    return devices;
//}
