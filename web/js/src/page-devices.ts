window.addEventListener("load", () => {
    console.warn("window onload event");
    updateApiDevices();
});

window.addEventListener("focus", () => {
    console.warn("window focus event");
    updateApiDevices();
});

function updateApiDevices() {
    console.warn((window as any).api as Api);
    // TODO: Fetch devices from GET "/api/devices"
}
