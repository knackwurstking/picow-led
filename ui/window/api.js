/**
 * @return {import("../types.d.ts").Api}
 */
export default () => {
    /**
     * @param {string} path
     * @returns {string}
     */
    function getUrl(path) {
        return `{{ .ServerPathPrefix }}${path}`;
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
        return _handleResponse(resp, url);
    }

    /**
     * @returns {Promise<import("../types.d.ts").ColorCache>}
     */
    async function color() {
        const url = getUrl("/api/color");
        const resp = await fetch(url);
        return _handleResponse(resp, url);
    }

    /**
     * @param {string} name
     * @returns {Promise<import("../types.d.ts").Color>}
     */
    async function colorForName(name) {
        const url = getUrl(`/api/color/${name}`);
        const resp = await fetch(url);
        return _handleResponse(resp, url);
    }

    // TODO: Add: POST "/api/color/:name" <- `number[]`

    return {
        devices,
        setDevicesColor,
        color,
        colorForName,
    };
};
