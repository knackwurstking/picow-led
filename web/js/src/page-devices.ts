window.addEventListener("load", () => {
    console.warn("window onload event");
    updateApiDevices();
});

window.addEventListener("focus", () => {
    console.warn("window focus event");
    updateApiDevices();
});

function updateApiDevices() {
    // FIXME: How to get rid of these error
    console.warn(api);
    // TODO: Fetch devices from GET "/api/devices"
}
