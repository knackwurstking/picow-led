import * as deviceItem from "./device-item";

window.onpageshow = async () => {
    setupAppBar();
    setupDevicesList();
};

async function setupAppBar() {
    const items = window.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"]!.innerText = "Devices";
}

async function setupDevicesList() {
    let devices: Device[] = [];

    const devicesList = document.querySelector<HTMLElement>(
        "._content.devices > .list",
    )!;

    const powerButtonToggle = async (device: Device) => {
        let color: Color;
        if (Math.max(...(device.color || [])) > 0) {
            color = (device.pins || device.color || []).map(() => 0);
        } else {
            color = currentColorForDevice(device);
        }

        await window.api.setDevicesColor(color, device);
    };

    window.store.obj.listen("devices", (data) => {
        devices = data;
        devicesList.innerHTML = "";

        devices.forEach((device) => {
            const item = deviceItem.create(device, () => {
                powerButtonToggle(device);
            });

            devicesList.appendChild(item);
        });
    });

    let timeout: NodeJS.Timeout | null = null;
    const onFocus = () => {
        if (timeout !== null) {
            clearTimeout(timeout);
            timeout = null;
        }

        timeout = setTimeout(async () => {
            await window.api.devices();

            timeout = null;
        });
    };
    window.onfocus = () => onFocus();
    onFocus();
}

function currentColorForDevice(device: Device): Color {
    return (
        window.store.currentDeviceColor(device.server.addr) ||
        (device.pins || []).map(() => 255)
    );
}
