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
        const status = resp.status;

        /** @type {any} */
        const respData = await resp.json();

        console.debug(`Got data from "${url}":`, respData);

        if ("message" in respData) {
            throw new Error(`${status}: ${respData.message}`);
        }

        return respData;
    }

    const api = {
        setDevicesColor,
    };

    // @ts-ignore
    window.api = api;
})();
//{{ end }}
