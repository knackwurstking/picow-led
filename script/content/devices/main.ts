document.addEventListener("DOMContentLoaded", async () => {
    setupAppBar();

    // TODO: Power Button:
    //  - query all device list item power buttons,
    //  - handle data-state
    //  - power button onclick event
});

async function setupAppBar() {
    const items = window.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"]!.innerText = "Devices";
}

//async function powerButtonToggle(device: Device) {
//    let color: Color;
//    if (Math.max(...(device.color || [])) > 0) {
//        color = (device.pins || device.color || []).map(() => 0);
//    } else {
//        color = currentColorForDevice(device);
//    }
//
//    await window.api.setDevicesColor(color, device);
//}

//function currentColorForDevice(device: Device): Color {
//    return (
//        window.store.currentDeviceColor(device.server.addr) ||
//        (device.pins || []).map(() => 255)
//    );
//}
