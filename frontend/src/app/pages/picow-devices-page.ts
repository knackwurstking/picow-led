import { customElement } from "lit/decorators.js";

import { css, html, PropertyValues, TemplateResult } from "lit";
import { CleanUp, Events, globalStylesToShadowRoot, UIStackLayoutPage } from "ui";

import * as app from "@app";
import * as lib from "@lib";
import * as types from "@types";

@customElement("picow-devices-page")
class PicowDevicesPage extends UIStackLayoutPage {
    name = "devices";

    // NOTE: For now the events object needs to be passed before the
    //       connectedCallback method is running
    public picowAppEvents: Events<types.AppBarEvents> | null = null;

    private store: types.PicowStore = document.querySelector(`ui-store`)!;
    private cleanup = new CleanUp();

    static get styles() {
        return css`
            ${UIStackLayoutPage.styles}

            :host {
                padding-top: var(--ui-app-bar-height);
                overflow: auto;
            }
        `;
    }

    protected render(): TemplateResult<1> {
        return html`<ul style="border-radius: var(--ui-radius);">
            <slot></slot>
        </ul>`;
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        super.firstUpdated(_changedProperties);
        globalStylesToShadowRoot(this.shadowRoot!);
    }

    connectedCallback(): void {
        super.connectedCallback();

        if (this.picowAppEvents !== null) {
            this.cleanup.add(
                this.picowAppEvents.addListener("add", async () => {
                    const dialog = new app.PicowDeviceSetupDialog();
                    dialog.allowDeletion = false;
                    dialog.open = true;
                    document.body.appendChild(dialog);

                    dialog.addEventListener("submit", async () => {
                        if (!dialog.device) return;
                        lib.ws.request("POST api.device", dialog.device);
                    });
                }),
            );
        }

        this.cleanup.add(
            // Store Events
            this.store.addListener("devices", (devices) => {
                const list = this.shadowRoot!.querySelector("ul")!;
                while (!!list.firstChild) list.removeChild(list.firstChild);
                for (const device of devices) {
                    setTimeout(() => {
                        const deviceItem = new app.PicowDeviceItem();
                        deviceItem.device = device;
                        list.appendChild(deviceItem);
                    });
                }
            }),

            // WS Events
            lib.ws.events.addListener("message-devices", async (data) => {
                this.store.setData("devices", data);
            }),
        );

        const getDevicesFromWS = async () => {
            try {
                await lib.ws.request("GET api.devices");
            } catch (err) {
                if (err instanceof Error) {
                    console.error(err);
                    lib.throwAlert({
                        message: err.message,
                        variant: "error",
                    });
                }
            }
        };

        getDevicesFromWS().then(() => {
            this.cleanup.add(lib.ws.events.addListener("open", getDevicesFromWS));
        });
    }

    disconnectedCallback(): void {
        super.disconnectedCallback();
        this.cleanup.run();
    }
}

export default PicowDevicesPage;
