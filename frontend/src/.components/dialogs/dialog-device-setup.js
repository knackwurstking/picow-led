import { html, UIDialog } from "ui";

/**
 * @extends {UIDialog<import("ui").UIDialog_Events & { "submit": Device, "delete": Device }>}
 */
export class DialogDeviceSetup extends UIDialog {
    static register = () => {
        customElements.define("dialog-device-setup", DialogDeviceSetup);
    };

    /**
     * @param {object} options
     * @param {string} [options.name]
     * @param {string} [options.addr]
     * @param {number[]} [options.pins]
     * @param {boolean} [options.allowDeletion]
     */
    constructor({
        name = "",
        addr = "",
        pins = [],
        allowDeletion = false,
    } = {}) {
        super("Device Setup");

        /** @type {Device} */
        this.data;
        /** @type {boolean} */
        this.allowDeletion;

        this.set({ name, addr, pins, allowDeletion });
    }

    /**
     * @param {object} options
     * @param {string} [options.name]
     * @param {string} [options.addr]
     * @param {number[]} [options.pins]
     * @param {boolean} [options.allowDeletion]
     */
    set({ name = "", addr = "", pins = [], allowDeletion = false } = {}) {
        this.data = {
            server: {
                name: name,
                addr: addr,
            },
            pins: pins,
        };

        this.allowDeletion = allowDeletion;
        this.render();
    }

    render() {
        this.innerHTML = html`
            <ui-flex-grid gap="0.5rem">
                <ui-flex-grid-item>
                    <ui-input
                        name="server.name"
                        type="text"
                        title="Server Name"
                        value="${this.data.server.name}"
                    ></ui-input>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-input
                        name="server.addr"
                        type="text"
                        title="Server Address"
                        value="${this.data.server.addr}"
                    ></ui-input>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-input
                        name="pins"
                        type="text"
                        title="GPIO pins in use"
                        placeholder="ex.: 0,1,2,3"
                        value="${this.data.pins.join(",")}"
                    ></ui-input>
                </ui-flex-grid-item>
            </ui-flex-grid>
        `;

        this.renderInputs();
        this.renderDeleteAction();
        this.renderCancelAction();
        this.renderSubmitAction();
    }

    /** @private */
    renderInputs() {
        const inputs = this.querySelectorAll(`ui-input`);
        inputs.forEach((/** @type {import("ui").UIInput} */ child) => {
            child.ui.events.on("input", (value) => {
                switch (child.getAttribute("name")) {
                    case "server.name":
                        this.data.server.name = value;
                        break;
                    case "server.addr":
                        this.data.server.addr = value;
                        break;
                    case "pins":
                        this.data.pins = value
                            .split(/,|\.|\s/)
                            .map((v) => parseInt(v))
                            .filter((v) => !isNaN(v));
                        break;
                }
            });
        });
    }

    /** @private */
    renderDeleteAction() {
        if (!this.allowDeletion) return;

        const el = UIDialog.createAction({
            variant: "full",
            color: "destructive",
            onClick: () => {
                this.ui.events.dispatch("delete", this.data);
                this.ui.close();
            },
        });
        el.action.innerHTML = "Delete";
        this.appendChild(el.container);
    }

    /** @private */
    renderCancelAction() {
        const cancel = UIDialog.createAction({
            variant: "full",
            color: "secondary",
            onClick: () => {
                this.ui.close();
            },
        });
        cancel.action.innerHTML = "Cancel";
        this.appendChild(cancel.container);
    }

    /** @private */
    renderSubmitAction() {
        const submit = UIDialog.createAction({
            variant: "full",
            color: "primary",
            onClick: () => {
                /** @type {import("ui").UIInput} */
                let addrInput = this.querySelector(
                    `ui-input[name="server.addr"]`
                );
                if (!this.data.server.addr) {
                    addrInput.ui.invalid = true;
                    return;
                }
                addrInput.ui.invalid = false;

                /** @type {import("ui").UIInput} */
                let pinsInput = this.querySelector(`ui-input[name="pins"]`);
                if (!this.data.pins.length) {
                    pinsInput.ui.invalid = true;
                    return;
                }
                pinsInput.ui.invalid = false;

                this.ui.events.dispatch("submit", this.data);
                this.ui.close();
            },
        });
        submit.action.innerHTML = "Submit";
        this.appendChild(submit.container);
    }
}
