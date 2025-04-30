//{{ define "script-window-api" }}
(() => {
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
     * @returns {Promise<Device[]>}
     */
    async function devices() {
        const url = getUrl("/api/devices");

        const resp = await fetch(url);
        return await _handleResponse(resp, url);
    }

    /**
     * @param {MicroColor | undefined | null} color
     * @param {Device[]} devices
     * @returns {Promise<Device[]>}
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

    /** @type {Api} */
    const api = {
        devices,
        setDevicesColor,
    };

    // @ts-ignore
    window.api = api;
})();
//{{ end }}
