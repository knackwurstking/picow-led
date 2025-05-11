import * as deviceItem from "./device-item";

function currentColorForDevice(device: Device): Color {
    return (
        window.store.currentDeviceColor(device.server.addr) ||
        (device.pins || []).map(() => 255)
    );
}

async function setupAppBar() {
    const items = window.utils.setupAppBarItems(
        "online-indicator",
        "title",
        "settings-button",
    );

    items["title"]!.innerText = "Devices";
}

async function setupDevicesList() {
    const devices = await window.api.devices();

    const devicesList = document.querySelector<HTMLElement>(
        "._content.devices > .list",
    );
    if (!devicesList) throw new Error(`devices list container is null`);

    devicesList.innerHTML = "";

    const onClick = async (device: Device) => {
        let color: Color;
        if (Math.max(...(device.color || [])) > 0) {
            color = (device.pins || device.color || []).map(() => 0);
        } else {
            color = currentColorForDevice(device);
        }

        await window.api.setDevicesColor(color, device);
    };

    devices.forEach((device) => {
        const item = deviceItem.create(device, () => {
            onClick(device);
        });

        devicesList.appendChild(item);
    });

    window.ws.events.addListener("device", (device) => {
        let child: HTMLElement;
        for (let x = 0; x < devices.length; x++) {
            if (devices[x].server.addr !== device.server.addr) {
                continue;
            }

            child = devicesList.children[x] as HTMLElement;

            deviceItem.update(child, device, () => {
                onClick(device);
            });
        }
    });
}

window.addEventListener("pageshow", async () => {
    setupAppBar();
    setupDevicesList();
});
