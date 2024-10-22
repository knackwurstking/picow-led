import { css, html, LitElement } from "lit";
import { customElement, property } from "lit/decorators.js";
import { svg } from "ui";
import { ws, WSEventsDevice } from "../../../lib/websocket";
import { PicowDeviceSetupDialog } from "../../dialogs/picow-device-setup-dialog";

/**
 * **Tag**: picow-options-button
 *
 * **Attributes**:
 *  - device: `WSEventsDevice` - [json]
 */
@customElement("picow-options-button")
export class PicowOptionsButton extends LitElement {
    @property({ type: Object, attribute: "device", reflect: true })
    device?: WSEventsDevice;

    static get styles() {
        return css`
            :host {
                height: 100%;
            }
        `;
    }

    protected render() {
        return html`
            <ui-icon-button
                ghost
                ripple
                @click=${async (ev: MouseEvent) => {
                    ev.stopPropagation();
                    if (!this.device) return;

                    const dialog = new PicowDeviceSetupDialog();

                    dialog.device = {
                        ...this.device,
                        server: { ...this.device.server },
                    };

                    dialog.allowDeletion = true;
                    dialog.open = true;
                    document.body.appendChild(dialog);

                    const validateDevice = () => {
                        if (!dialog.device) {
                            throw new Error(
                                `missing dialog data: device undefined`,
                            );
                        }
                    };

                    dialog.addEventListener("delete", async () => {
                        validateDevice();

                        ws.request("DELETE api.device", {
                            addr: dialog.device!.server.addr,
                        });
                    });

                    dialog.addEventListener("submit", async () => {
                        validateDevice();
                        ws.request("PUT api.device", dialog.device!);
                    });
                }}
            >
                ${svg.smoothieLineIcons.moreVertical}
            </ui-icon-button>
        `;
    }
}
