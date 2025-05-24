document.addEventListener("DOMContentLoaded", async () => {
    setupAppBar();

    const deviceListItems = document.querySelectorAll<HTMLButtonElement>(
        `.device-list-item button.power`,
    );

    deviceListItems.forEach((button) => {
        button.onclick = async () => {
            const addr = button.getAttribute(`data-addr`)!;
            const pins = JSON.parse(button.getAttribute(`data-pins`)!);
            let color = JSON.parse(button.getAttribute("data-color")!);

            if (color && Math.max(...color) > 0) {
                color = (pins || []).map(() => 0);
            } else {
                color = currentColor(addr, pins || []);
            }

            await window.api.setDevicesColor(color, addr);
        };
    });
});

async function setupAppBar() {
    const items = window.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"]!.innerText = "Devices";
}

function currentColor(addr: string, pins: number[]): Color {
    return window.store.currentDeviceColor(addr) || pins.map(() => 255);
}
