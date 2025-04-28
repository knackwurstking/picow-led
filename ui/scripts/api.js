(() => {
    function getUrl() {
        return ``; // TODO: Add the server prefix somehow, or just move this crap back to templ
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

        const url = getUrl() + "/api/devices/color";
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
        let respData;
        try {
            respData = await resp.json();
        } catch (err) {
            console.error(err); // TODO: Error handling, (notifications)
        }

        console.debug(`Got data from "${url}":`, respData);

        if ("message" in respData) {
            console.error(`${status}: ${respData.message}`); // TODO: Error handling, (notifications)
            return;
        }

        return respData;
    }

    /** @type {Api} */
    const api = {
        setDevicesColor,
    };

    // @ts-ignore
    window.api = api;
})();
