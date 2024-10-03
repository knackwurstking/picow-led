import { moreVertical as svgOptions } from "ui/svg/smoothie-line-icons";

import { html, UIIconButton } from "ui";
import * as api from "../../../lib/api";
import createDeviceSetupDialog from "../../dialogs/createDeviceSetupDialog";

export interface PicowOptionsButton_Picow {
    root: PicowOptionsButton;
    set(device: Device): void;
}

export default class PicowOptionsButton extends UIIconButton {
    device: Device;
    store: PicowStore;
    picow: PicowOptionsButton_Picow;

    constructor(device: Device | null = null) {
        super();

        this.device = device;
        this.store = document.querySelector(`ui-store`);

        this.picow = {
            root: this,

            set(device) {
                this.root.device = device;
            },
        };

        this.#render();
    }

    #render() {
        this.ui.ghost = true;

        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    height: 100%;
                }
            </style>
        `;

        this.innerHTML = svgOptions;

        this.onclick = async (ev) => {
            ev.stopPropagation();

            const setupDialog = await createDeviceSetupDialog({
                name: this.device.server.name,
                addr: this.device.server.addr,
                pins: this.device.pins,
                allowDeletion: true,
            });

            setupDialog.events.on("delete", async (deviceToDelete) => {
                await api.Delete(this.store, "/api/device", deviceToDelete);
            });

            setupDialog.events.on("submit", async (deviceToSubmit) => {
                await api.Put(this.store, "/api/device", deviceToSubmit);
            });

            setupDialog.open();
        };
    }
}

console.debug(`Register the "picow-options-button"`);
customElements.define("picow-options-button", PicowOptionsButton);
