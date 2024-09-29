import { Events, html, UIDialog } from "ui";

/**
 * @typedef DeviceSetupDialog_Events
 * @type {import("ui").UIDialog_Events & {
 *  "submit": Device;
 *  "delete": Device;
 *  }}
 *
 * @typedef DeviceSetupDialog_Options
 * @type {{
 *  name?: string;
 *  addr?: string;
 *  pins?: number[];
 *  allowDeletion?: boolean;
 * }}
 *
 * @typedef {{
 *  element: UIDialog;
 *  get events(): Events<DeviceSetupDialog_Events>;
 *  set(options: DeviceSetupDialog_Options): void;
 *  open(): void;
 *  close(): void;
 * }} DeviceSetupDialog
 */

/**
 * @param {DeviceSetupDialog_Options} options
 * @returns {Promise<DeviceSetupDialog>}
 */
export default async function ({
    name = "",
    addr = "",
    pins = [],
    allowDeletion = false,
}) {
    /**
     * @type {UIDialog<DeviceSetupDialog_Events>}
     */
    const dialog = new UIDialog("Device Setup");

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
            const inputs = dialog.querySelectorAll(`ui-input`);
            inputs.forEach((/** @type {import("ui").UIInput} */ input) => {
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

        {
            if (!allowDeletion) return;

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
                    /** @type {import("ui").UIInput} */
                    let addrInput = dialog.querySelector(
                        `ui-input[name="server.addr"]`
                    );
                    if (!device.server.addr) {
                        addrInput.ui.invalid = true;
                        return;
                    }
                    addrInput.ui.invalid = false;

                    /** @type {import("ui").UIInput} */
                    let pinsInput = dialog.querySelector(
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

    /**
     * @param {DeviceSetupDialog_Options} options
     */
    function set({ name = "", addr = "", pins = [], allowDeletion = false }) {
        device.server.name = name;
        device.server.addr = addr;
        device.pins = pins;
        allowDeletion = allowDeletion;
        render();
    }

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
