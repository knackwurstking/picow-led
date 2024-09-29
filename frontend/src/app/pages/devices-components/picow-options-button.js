import { moreVertical as svgOptions } from "ui/svg/smoothie-line-icons";
import { html, UIIconButton } from "ui";
import createDeviceSetupDialog from "../../dialogs/createDeviceSetupDialog";
import { utils, api } from "../../../lib";

export default class PicowOptionsButton extends UIIconButton {
    /**
     * @param {Device | null} [device]
     */
    constructor(device = null) {
        super();

        /**
         * @type {Device | null}
         */
        this.device = device;

        /**
         * @type {PicowStore}
         */
        this.store = document.querySelector(`ui-store`);

        this.picow = {
            root: this,

            /**
             * @param {Device} device
             */
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
