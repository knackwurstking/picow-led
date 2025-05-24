document.addEventListener("DOMContentLoaded", async () => {
    setupAppBar();

    const powerButtons = document.querySelectorAll<HTMLButtonElement>(
        `.device-list-item button.power`,
    );

    powerButtons.forEach((button) => {
        {
            const color = JSON.parse(button.getAttribute("data-color")!);
            setPowerButtonState(button, color);
            setPowerButtonBackground(button, color);
        }

        button.onclick = async () => {
            if (button.getAttribute(`data-state`) === "processing") return;
            button.setAttribute("data-state", "processing");

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

    window.ws.events.addListener("message", async (data) => {
        switch (data.type) {
            case "device":
                updatePowerButton(data.data);
                break;
        }
    });
});

function setupAppBar() {
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

function updatePowerButton(device: Device) {
    const powerButtons = document.querySelectorAll<HTMLButtonElement>(
        `.device-list-item button.power`,
    );

    powerButtons.forEach((button) => {
        if (button.getAttribute("data-addr") !== device.server.addr) return;

        button.setAttribute(`data-color`, JSON.stringify(device.color));

        button.setAttribute(`data-pins`, JSON.stringify(device.pins));

        setPowerButtonState(button, device.color);
        setPowerButtonBackground(button, device.color);
    });
}

function setPowerButtonState(
    button: HTMLButtonElement,
    color: number[] | null,
) {
    button.setAttribute(
        `data-state`,
        Math.max(...(color || [])) > 0 ? "on" : "off",
    );
}

function setPowerButtonBackground(
    button: HTMLButtonElement,
    color: number[] | null,
) {
    button.querySelector<HTMLElement>(`.background`)!.style.backgroundColor =
        `rgb(${(color || [0, 0, 0]).slice(0, 3).join(", ")})`;
}
