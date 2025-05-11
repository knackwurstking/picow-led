type onClickPowerButton = (item: HTMLElement) => void | Promise<void>;

export function create(
    device: Device,
    onClickPowerButton: onClickPowerButton,
): HTMLElement {
    const t = document.querySelector<HTMLTemplateElement>(
        `template[name="device-list-item"]`,
    );
    if (!t) {
        throw new Error(
            `Nope, template with name "color-range-slider" not found`,
        );
    }

    const item = (
        t.content.cloneNode(true) as HTMLElement
    ).querySelector<HTMLElement>(".device-list-item");
    if (!item) throw new Error(`device list item is null`);

    return update(item, device, onClickPowerButton);
}

export function update(
    item: HTMLElement,
    device: Device,
    onClickPowerButton: onClickPowerButton,
): HTMLElement {
    item.setAttribute("data-addr", device.server.addr);

    const title = item.querySelector<HTMLElement>(`.title`);
    if (!title) throw new Error(`title element is null`);

    title.innerHTML = device.server.name || device.server.addr;

    const editButton = item.querySelector<HTMLElement>(`button.edit`);
    if (!editButton) throw new Error(`edit button element is null`);

    editButton.setAttribute("data-addr", device.server.addr);

    const powerButton =
        item.querySelector<HTMLButtonElement>(`button.power-button`);
    if (!powerButton) throw new Error(`power button element is null`);

    powerButton.onclick = async (ev) => {
        const target = ev.currentTarget as HTMLButtonElement;

        if (target.getAttribute("data-state") === "processing") {
            return;
        }

        target.setAttribute("data-state", "processing");

        await onClickPowerButton(item);
    };

    const background = powerButton.querySelector<HTMLElement>(`.background`);
    if (!background) throw new Error(`background element is null`);

    background.style.backgroundColor = `rgb(${[...(device.color || []), 0, 0, 0].slice(0, 3).join(", ")})`;

    if (Math.max(...(device.color || []), 0, 0, 0)) {
        powerButton.setAttribute("data-state", "on");
    } else {
        powerButton.setAttribute("data-state", "off");
    }

    return item;
}
