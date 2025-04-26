// TODO: Fetch devices data from the api, on focus and on window load

window.addEventListener("load", () => {
    console.debug("window onload event");
    updateApiDevices();
});

window.addEventListener("focus", () => {
    console.debug("window focus event");
    updateApiDevices();
});

function updateApiDevices() {
    // ...
}
