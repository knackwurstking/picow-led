function getUrl() {
    return ``; // TODO: Add the server prefix somehow, or just move this crap back to templ
}

/**
 * @param {MicroColor | undefined | null} color
 * @param {Device[]} devices
 * @returns {Promise<void>}
 */
async function color(color, ...devices) {
    if (!color) {
        color = [255, 255, 255, 255];
    }

    const url = getUrl() + "/api/devices/color";
    const data = { devices, color };
    console.debug(`POST "${url}":`, data);

    /** @type {number} */
    let status;
    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    })
        .then((resp) => {
            status = resp.status;
            return resp.json();
        })
        .then((data) => {
            console.debug(`Got data from "${url}":`, data);

            if ("message" in data) {
                throw new Error(`${status}: ${data.message}`);
            }

            console.warn(data); // TODO: Response handling
        })
        .catch((err) => {
            console.error(err); // TODO: Error handling, (notifications)
        });
}

// @ts-ignore
window.api = {
    color,
};
