import { html, UIDialog, UIInput, type Events, type UIDialog_Events } from "ui";

export type DeviceSetupDialog_Events = {
    submit: Device;
    delete: Device;
} & UIDialog_Events;

export interface DeviceSetupDialog_Options {
    name?: string;
    addr?: string;
    pins?: number[];
    allowDeletion?: boolean;
}

export interface DeviceSetupDialog {
    element: UIDialog;
    get events(): Events<DeviceSetupDialog_Events>;
    set(options: DeviceSetupDialog_Options): void;
    open(): void;
    close(): void;
}

export default async function ({
    name = "",
    addr = "",
    pins = [],
    allowDeletion = false,
}: DeviceSetupDialog_Options): Promise<DeviceSetupDialog> {
    const dialog = new UIDialog<DeviceSetupDialog_Events>("Device Setup");

    const device = {
        server: {
            name: name,
            addr: addr,
        },
        pins: pins,
    };

    function render() {
        dialog.innerHTML = html`
            <ui-flex-grid gap="0.5rem">
                <ui-flex-grid-item>
                    <ui-input
                        name="server.name"
                        type="text"
                        title="Server Name"
                        value="${device.server.name}"
                    ></ui-input>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-input
                        name="server.addr"
                        type="text"
                        title="Server Address"
                        value="${device.server.addr}"
                    ></ui-input>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-input
                        name="pins"
                        type="text"
                        title="GPIO pins in use"
                        placeholder="ex.: 0,1,2,3"
                        value="${device.pins.join(",")}"
                    ></ui-input>
                </ui-flex-grid-item>
            </ui-flex-grid>
        `;

        // Clean up all actions first (if set was run twice)
        dialog
            .querySelectorAll(`[slot="actions"]`)
            .forEach((el) => dialog.removeChild(el));

        // --------------------- //
        // Render Input Elements //
        // --------------------- //

        {
            const inputs = dialog.querySelectorAll<UIInput>(`ui-input`);
            inputs.forEach((input) => {
                input.ui.events.on("input", (value) => {
                    switch (input.getAttribute("name")) {
                        case "server.name":
                            device.server.name = value;
                            break;

                        case "server.addr":
                            device.server.addr = value;
                            break;

                        case "pins":
                            device.pins = value
                                .split(/,|\.|\s/)
                                .map((v) => parseInt(v))
                                .filter((v) => !isNaN(v));
                            break;
                    }
                });
            });
        }

        // -------------------- //
        // Render Delete Action //
        // --------------------- //

        if (allowDeletion) {
            const el = UIDialog.createAction({
                variant: "full",
                color: "destructive",
                onClick: () => {
                    dialog.ui.events.dispatch("delete", device);
                    dialog.ui.close();
                },
            });

            el.action.innerHTML = "Delete";
            dialog.appendChild(el.container);
        }

        // -------------------- //
        // Render Cancel Action //
        // --------------------- //

        {
            const el = UIDialog.createAction({
                variant: "full",
                color: "secondary",
                onClick: () => {
                    dialog.ui.close();
                },
            });
            el.action.innerHTML = "Cancel";
            dialog.appendChild(el.container);
        }

        // -------------------- //
        // Render Submit Action //
        // --------------------- //

        {
            const el = UIDialog.createAction({
                variant: "full",
                color: "primary",
                onClick: () => {
                    let addrInput = dialog.querySelector<UIInput>(
                        `ui-input[name="server.addr"]`
                    );

                    if (!device.server.addr) {
                        addrInput.ui.invalid = true;
                        return;
                    }

                    addrInput.ui.invalid = false;

                    let pinsInput = dialog.querySelector<UIInput>(
                        `ui-input[name="pins"]`
                    );

                    if (!device.pins.length) {
                        pinsInput.ui.invalid = true;
                        return;
                    }

                    pinsInput.ui.invalid = false;

                    dialog.ui.events.dispatch("submit", device);
                    dialog.ui.close();
                },
            });
            el.action.innerHTML = "Submit";
            dialog.appendChild(el.container);
        }
    }

    function set({
        name = "",
        addr = "",
        pins = [],
        allowDeletion = false,
    }: DeviceSetupDialog_Options) {
        device.server.name = name;
        device.server.addr = addr;
        device.pins = pins;
        allowDeletion = allowDeletion;
        render();
    }

    set({ name, addr, pins, allowDeletion });
    document.body.appendChild(dialog);
    dialog.ui.events.on("close", () => {
        document.body.removeChild(dialog);
    });

    return {
        element: dialog,

        get events() {
            return dialog.ui.events;
        },

        set,

        open() {
            dialog.ui.open(true);
        },

        close() {
            dialog.ui.close();
        },
    };
}
