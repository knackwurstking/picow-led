import { html, LitElement, PropertyValues } from "lit";
import { customElement, property } from "lit/decorators.js";
import { UIDialog, UIInput } from "ui";
import { WSEventsDevice } from "../../lib/websocket";

/**
 * **Tag**: picow-device-setup-dialog
 *
 * **Attributes**:
 *  - device: `WSEventsDevice` - [json]
 *  - allow-deletion: `boolean`
 *  - open: `boolean`
 *
 * **Events**:
 *  - submit
 *  - delete
 *
 * **Public Methods**:
 *  - `rootElement(): UIDialog`
 *  - `async open()`
 *  - `async close()`
 */
@customElement("picow-device-setup-dialog")
export class PicowDeviceSetupDialog extends LitElement {
    @property({ type: Object, attribute: "device", reflect: true })
    device?: WSEventsDevice;

    @property({ type: Boolean, attribute: "allow-deletion", reflect: true })
    allowDeletion: boolean = false;

    @property({ type: Boolean, attribute: "open", reflect: true })
    open: boolean = false;

    protected render() {
        return html`
            <ui-dialog
                title="Device Setup"
                ?open=${this.open}
                modal
                inert
                @close=${async () => {
                    try {
                        document.removeChild(this);
                    } catch {}
                }}
            >
                <ui-flex-grid gap="0.5rem">
                    <ui-flex-grid-item>
                        <ui-input
                            name="server.name"
                            type="text"
                            title="Server Name"
                            value="${this.device?.server.name}"
                            @input=${async (ev: Event) => {
                                if (!this.device) return;
                                const target = ev.currentTarget as UIInput;
                                this.device.server.name = target.value;
                            }}
                        ></ui-input>
                    </ui-flex-grid-item>

                    <ui-flex-grid-item>
                        <ui-input
                            name="server.addr"
                            type="text"
                            title="Server Address"
                            value="${this.device?.server.addr}"
                            @input=${async (ev: Event) => {
                                if (!this.device) return;
                                const target = ev.currentTarget as UIInput;
                                this.device.server.addr = target.value;
                            }}
                        ></ui-input>
                    </ui-flex-grid-item>

                    <ui-flex-grid-item>
                        <ui-input
                            name="pins"
                            type="text"
                            title="GPIO pins in use"
                            placeholder="ex.: 0,1,2,3"
                            value="${this.device?.pins?.join(",") || ""}"
                            @input=${async (ev: Event) => {
                                if (!this.device) return;
                                const target = ev.currentTarget as UIInput;
                                this.device.pins = target.value
                                    .split(/,|\.|\s/)
                                    .map((v) => parseInt(v))
                                    .filter((v) => !isNaN(v));
                            }}
                        ></ui-input>
                    </ui-flex-grid-item>
                </ui-flex-grid>
            </ui-dialog>
        `;
    }

    protected updated(_changedProperties: PropertyValues): void {
        const rootElement = this.rootElement();

        if (this.allowDeletion) {
            // Create "Delete" action
            rootElement.addDialogActionButton("Delete", {
                onClick: async () => {
                    this.dispatchEvent(new Event("delete"));
                    rootElement.close();
                },
                variant: "full",
                color: "destructive",
                flex: 0,
            });
        }

        // Create "Cancel" action
        rootElement.addDialogActionButton("Cancel", {
            onClick: async () => {
                rootElement.close();
            },
            variant: "full",
            color: "secondary",
            flex: 0,
        });

        // Create "Submit" action
        rootElement.addDialogActionButton("Submit", {
            onClick: async () => {
                let addrInput = this.shadowRoot!.querySelector<UIInput>(
                    `ui-input[name="server.addr"]`,
                )!;

                if (!this.device?.server.addr) {
                    addrInput.invalid = true;
                    return;
                }

                addrInput.invalid = false;

                let pinsInput = this.shadowRoot!.querySelector<UIInput>(
                    `ui-input[name="pins"]`,
                )!;

                if (!this.device?.pins?.length) {
                    pinsInput.invalid = true;
                    return;
                }

                pinsInput.invalid = false;

                this.dispatchEvent(new Event("submit"));
                rootElement.close();
            },
            variant: "full",
            color: "primary",
            flex: 0,
        });
    }

    public rootElement(): UIDialog {
        return this.shadowRoot!.querySelector(`ui-dialog`)!;
    }

    public async show() {
        if (
            [...document.body.children].find((child) => child === this) ===
            undefined
        ) {
            document.body.appendChild(this);
        }

        const rootElement = this.rootElement();
        if (rootElement === null) {
            this.open = true;
            return;
        }

        rootElement.show({ modal: true, inert: true });
    }

    public async close() {
        this.rootElement().close();
    }
}
