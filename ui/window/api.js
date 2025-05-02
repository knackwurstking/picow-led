/**
 * @returns {import("../types.d.ts").Api}
 */
export function create() {
    /** @type {import("../types.d.ts").PageWindow} */
    // @ts-ignore
    const w = window;

    /**
     * @param {string} path
     * @returns {string}
     */
    function getUrl(path) {
        // @ts-ignore
        return process.env.SERVER_PATH_PREFIX + `${path}`;
    }

    /**
     * @param {Response} resp
     * @param {string} url
     * @returns {Promise<any>}
     */
    async function _handleResponse(resp, url) {
        const status = resp.status;

        if (!resp.ok) {
            throw new Error(`${status}: ${(await resp.text()) || "???"}`);
        }

        const respData = await resp.json();
        console.debug(`Got data from "${url}":`, respData);
        return respData;
    }

    /**
     * @returns {Promise<import("../types.d.ts").Device[]>}
     */
    async function devices() {
        const url = getUrl("/api/devices");

        const resp = await fetch(url);
        return await _handleResponse(resp, url);
    }

    /**
     * @param {import("../types.d.ts").Color | undefined | null} color
     * @param {import("../types.d.ts").Device[]} devices
     * @returns {Promise<import("../types.d.ts").Device[]>}
     */
    async function setDevicesColor(color, ...devices) {
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

        /** @type {import("../types.d.ts").Device[]} */
        devices = await _handleResponse(resp, url);

        w.store.obj.update("devices", (d) => {
            for (let sI = 0; sI < d.length; sI++) {
                for (let i = 0; i < devices.length; i++) {
                    if (d[sI].server.addr === devices[i].server.addr) {
                        d[sI] = devices[i];
                    }
                }
            }

            return d;
        });

        return devices;
    }

    /**
     * @returns {Promise<import("../types.d.ts").ColorCache>}
     */
    async function colors() {
        const url = getUrl("/api/color");
        const resp = await fetch(url);
        return _handleResponse(resp, url);
    }

    /**
     * @param {number} index
     * @returns {Promise<import("../types.d.ts").Color>}
     */
    async function color(index) {
        const url = getUrl(`/api/color/${index}`);
        const resp = await fetch(url);
        return _handleResponse(resp, url);
    }

    /**
     * @param {number} index
     * @param {import("../types.d.ts").Color} color
     * @returns {Promise<void>}
     */
    async function setColor(index, color) {
        const url = getUrl(`/api/color/${index}`);

        const resp = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(color),
        });

        return _handleResponse(resp, url);
    }

    return {
        devices,
        setDevicesColor,
        colors,
        color,
        setColor,
    };
}
