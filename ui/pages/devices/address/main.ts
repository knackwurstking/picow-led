const colorStorageItem = require("./color-storage-item.js");
const colorRangeSlider = require("./color-range-slider.js");

function pageDeviceAddress(): string {
    return decodeURIComponent(location.pathname.split("/").reverse()[0]);
}

function pageDevice(): Device {
    const addr = pageDeviceAddress();
    const device = window.store.device(addr);
    if (!device) throw new Error(`device not found for ${addr}`);
    return device;
}

function pageCurrentColor(): Color {
    return (
        window.store.currentDeviceColor(pageDeviceAddress()) ||
        (pageDevice().pins || []).map(() => 255)
    );
}

function pagePickedColor(): Color | null {
    // Get color from active item
    const color = pageActiveColor().slice(0, 3);

    // Get range slider values
    color.push(...pageRangeSliderValues());

    return color;
}

function pageActiveColor(): Color {
    let color = [];
    const activeItem = document.querySelector(`.color-storage-item.active`);
    if (activeItem) {
        color.push(
            ...colorStorageItem.splitDataColor(
                activeItem.getAttribute("data-color"),
            ),
        );
    } else {
        color = [255, 255, 255];
    }
    return color;
}

function pageRangeSliderValues(): number[] {
    return Array.from(
        document.querySelectorAll<HTMLInputElement>(
            `.range-sliders .color-range-slider input[type="range"]`,
        ),
    ).map((input) => {
        return parseInt(input.value || "0", 10);
    });
}

function setupAppBar(): void {
    const device = pageDevice();
    const items = window.utils.setupAppBarItems("online-indicator", "title");
    items.title!.innerText = device ? device.server.name : "";
}

async function setupColorStorage(): Promise<void> {
    const colorStorageContainer = document.querySelector<HTMLElement>(
        `.color-storage-container`,
    );
    if (!colorStorageContainer)
        throw new Error(`color storage container is null`);

    colorStorageContainer.innerHTML = "";

    const currentColor = pageCurrentColor();
    const currentColorString = currentColor
        .slice(0, 3)
        .join(colorStorageItem.colorSeparator);

    // Create color storage items
    const device = pageDevice();
    (await window.api.colors()).forEach((color, index) => {
        const item = colorStorageItem.create(index, color, {
            device,

            onClick(color): void {
                const colorString = color.join(colorStorageItem.colorSeparator);

                Array.from(colorStorageContainer.children).forEach((child) => {
                    if (child.getAttribute("data-color") === colorString) {
                        if (!child.classList.contains("active")) {
                            child.classList.add("active");

                            const color = colorStorageItem.splitDataColor(
                                child.getAttribute("data-color"),
                            );

                            window.api.setDevicesColor(
                                [...color, ...pageRangeSliderValues()],
                                device,
                            );
                        }
                    } else {
                        child.classList.remove("active");
                    }
                });
            },

            onChange(color): void {
                if (!item.classList.contains("active")) {
                    return;
                }

                window.api.setDevicesColor(color, device);
            },
        });

        if (item.getAttribute("data-color") === currentColorString) {
            item.classList.add("active");
        }

        colorStorageContainer.appendChild(item);
    });
}

async function setupRangeSliders(): Promise<void> {
    const container = document.querySelector<HTMLElement>(".range-sliders");
    if (!container) throw new Error(`range sliders container is null`);

    container.innerHTML = "";

    const device = pageDevice();

    if (device.pins.length > 3) {
        container.style.display = "block";
    } else {
        container.style.display = "none";
        return;
    }

    if (device.pins) {
        const currentColor = pageCurrentColor();
        let timeout: NodeJS.Timeout | null = null;
        device.pins.slice(3).forEach((pin, index) => {
            index += 3;
            const slider = colorRangeSlider.create(
                `Pin: ${pin.toString()}`,
                currentColor[index] || 0,
                () => {
                    // NOTE: Update device color (api) with some timeout
                    //       (250ms?), i should use websockets for this later
                    if (timeout !== null) {
                        clearTimeout(timeout);
                        timeout = null;
                    }
                    timeout = setTimeout(() => {
                        timeout = null;
                        window.api.setDevicesColor(pagePickedColor(), device);
                    }, 250);
                },
            );

            container.appendChild(slider);
        });
    }
}

window.addEventListener("pageshow", async () => {
    setupAppBar();
    setupColorStorage();
    setupRangeSliders();

    console.debug("device address:", pageDevice());
    console.table({
        activeColor: pageActiveColor(),
        currentColor: pageCurrentColor(),
        pickedColor: pagePickedColor(),
        rangeSliderValues: pageRangeSliderValues(),
    });
});
