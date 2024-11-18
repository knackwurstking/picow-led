import { customElement, property } from "lit/decorators.js";

import { css, html, LitElement } from "lit";
import { svg } from "ui";

import * as app from "@app";
import * as lib from "@lib";
import * as types from "@types";

@customElement("picow-options-button")
class PicowOptionsButton extends LitElement {
    @property({ type: Object, attribute: "device", reflect: true })
    device?: types.WSEventsDevice;

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

                    const dialog = new app.PicowDeviceSetupDialog();

                    dialog.device = {
                        ...this.device,
                        server: { ...this.device.server },
                    };

                    dialog.allowDeletion = true;
                    dialog.open = true;
                    document.body.appendChild(dialog);

                    const validateDevice = () => {
                        if (!dialog.device) {
                            throw new Error(`missing dialog data: device undefined`);
                        }
                    };

                    dialog.addEventListener("delete", async () => {
                        validateDevice();

                        lib.ws.request("DELETE api.device", {
                            addr: dialog.device!.server.addr,
                        });
                    });

                    dialog.addEventListener("submit", async () => {
                        validateDevice();
                        lib.ws.request("PUT api.device", dialog.device!);
                    });
                }}
            >
                ${svg.smoothieLineIcons.moreVertical}
            </ui-icon-button>
        `;
    }
}

export default PicowOptionsButton;
