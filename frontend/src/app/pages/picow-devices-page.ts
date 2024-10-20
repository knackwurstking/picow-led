import { css, html, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import { CleanUp, UIStackLayoutPage } from "ui";
import { throwAlert } from "../../lib/utils";
import { ws } from "../../lib/websocket";
import { PicowStore } from "../../types";
import { PicowDeviceSetupDialog } from "../dialogs/picow-device-setup-dialog";
import { PicowAppBar } from "../picow-app-bar";
import { PicowDeviceItem } from "./devices-components/picow-device-item";

/**
 * **Tag**: picow-devices-page
 */
@customElement("picow-devices-page")
export class PicowDevicesPage extends UIStackLayoutPage {
    name = "devices";

    private appBar: PicowAppBar = document.querySelector(`ui-app-bar`)!;
    private store: PicowStore = document.querySelector(`ui-store`)!;
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

    connectedCallback(): void {
        super.connectedCallback();
        this.cleanup.add(
            // AppBar Events
            this.appBar.events.addListener("add", async () => {
                const dialog = new PicowDeviceSetupDialog();
                dialog.allowDeletion = false;

                dialog.addEventListener("submit", async () => {
                    if (!dialog.device) return;
                    ws.request("POST api.device", dialog.device);
                });

                dialog.open();
            }),

            // Store Events
            this.store.addListener("devices", (devices) => {
                const list = this.shadowRoot!.querySelector("ul")!;
                while (!!list.firstChild) list.removeChild(list.firstChild);
                for (const device of devices) {
                    setTimeout(() => {
                        const deviceItem = new PicowDeviceItem();
                        deviceItem.device = device;
                        list.appendChild(deviceItem);
                    });
                }
            }),

            // WS Events
            ws.events.addListener("message-devices", async (data) => {
                this.store.setData("devices", data);
            })
        );

        const getDevicesFromWS = async () => {
            try {
                await ws.request("GET api.devices");
            } catch (err) {
                if (err instanceof Error) {
                    console.error(err);
                    throwAlert({
                        message: err.message,
                        variant: "error",
                    });
                }
            }
        };

        getDevicesFromWS().then(() => {
            this.cleanup.add(ws.events.addListener("open", getDevicesFromWS));
        });
    }

    disconnectedCallback(): void {
        super.disconnectedCallback();
        this.cleanup.run();
    }
}
