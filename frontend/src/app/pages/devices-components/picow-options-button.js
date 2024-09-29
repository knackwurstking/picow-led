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
                const url = await api.url(this.store, "/api/device");

                try {
                    const resp = await fetch(url, {
                        method: "DELETE",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(deviceToDelete),
                    });

                    if (!resp.ok) {
                        resp.text().then((e) => {
                            const message = `Server response to ${url}: ${e}`;
                            utils.throwAlert({ message, variant: "error" });
                            console.error(message);
                        });

                        const message = `Fetch from "${url}" with status code ${resp.status}`;
                        console.error(message);
                        utils.throwAlert({ message, variant: "error" });
                    }
                } catch (err) {
                    utils.throwAlert({ message: err, variant: "error" });
                }
            });

            setupDialog.events.on("submit", async (deviceToSubmit) => {
                const url = await api.url(this.store, "/api/device");

                try {
                    const resp = await fetch(url, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(deviceToSubmit),
                    });

                    if (!resp.ok) {
                        resp.text().then((e) => {
                            const message = `Server response to ${url}: ${e}`;
                            utils.throwAlert({ message, variant: "error" });
                            console.error(message);
                        });

                        const message = `Fetch from "${url}" with status code ${resp.status}`;
                        console.error(message);
                        utils.throwAlert({ message, variant: "error" });
                    }
                } catch (err) {
                    utils.throwAlert({ message: err, variant: "error" });
                }
            });

            setupDialog.open();
        };
    }
}

console.debug(`Register the "picow-options-button"`);
customElements.define("picow-options-button", PicowOptionsButton);
