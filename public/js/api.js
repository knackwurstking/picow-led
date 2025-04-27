function url() {
    return ``; // TODO: ...
}

//function devices() {
//    const url = this.url() + "/api/devices";
//    // TODO: GET "/api/devices"
//    return [];
//}

function color(color, ...devices) {
    if (!color) {
        color = [255, 255, 255, 255];
    }

    const url = this.url() + "/api/devices/color";
    const data = { devices, color };
    console.debug(`POST "${url}":`, data)

    let status;
    fetch(url, {
        method: "POST",
        header: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data),
    }).then((resp) => {
        status = resp.status
        return resp.json()
    }).then((data) => {
        console.debug(`fetch "${url}":`, data)

        if ("message" in data) {
            throw new Error(`${status}: ${data.message}`)
        }

        console.warn(data) // TODO: Response handling
    }).catch((err) => {
        console.error(err) // TODO: Error handling, (notifications)
    })
}

window.api = {
    url,
    //devices,
    color,
}
