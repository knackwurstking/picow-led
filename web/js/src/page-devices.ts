window.addEventListener("load", () => {
    updateApiDevices();
});

window.addEventListener("focus", () => {
    updateApiDevices();
});

function updateApiDevices() {
    // TODO: Fetch devices from GET "/api/devices"
    console.warn((window as any).api as Api);
}
