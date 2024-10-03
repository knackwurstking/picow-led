import { CleanUp, html, UIStackLayoutPage } from "ui";
import { api, utils } from "../../lib";
import ws from "../../lib/websocket";
import createDeviceSetupDialog from "../dialogs/createDeviceSetupDialog";
import PicowDeviceItem from "./devices-components/picow-device-item";

export default class PicowDevicesPage extends UIStackLayoutPage {
    /**
     * @param {object} options
     * @param {import("../create-app-bar").AppBar} options.appBar
     */
    constructor({ appBar }) {
        super("devices");

        /**
         * @type {PicowStore}
         */
        this.store = document.querySelector(`ui-store`);
        this.appBar = appBar;
        this.cleanup = new CleanUp();

        this.#render();
    }

    #render() {
        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    padding-top: var(--ui-app-bar-height);
                    overflow: auto;
                }
            </style>
        `;

        this.innerHTML = html`
            <ul style="border-radius: var(--ui-radius);"></ul>
        `;
    }

    connectedCallback() {
        super.connectedCallback();

        this.cleanup.add(
            // -------------------- //
            // Handle AppBar events //
            // -------------------- //

            this.appBar.events.on("add", async () => {
                const setupDialog = await createDeviceSetupDialog({
                    allowDeletion: false,
                });

                setupDialog.events.on("submit", async (deviceToSubmit) => {
                    await api.Post(this.store, "/api/device", deviceToSubmit);
                });

                setupDialog.open();
            }),

            // ------------------- //
            // Handle Store events //
            // ------------------- //

            this.store.ui.on("devices", (devices) => {
                const list = this.querySelector("ul");
                while (!!list.firstChild) list.removeChild(list.firstChild);
                for (const device of devices) {
                    setTimeout(() => {
                        list.appendChild(new PicowDeviceItem(device));
                    });
                }
            }),

            // ----------------------- //
            // Handle WebSocket events //
            // ----------------------- //

            ws.events.on("open", async () => {
                try {
                    const devices = await this.fetchApiDevices();
                    if (!!devices) this.store.ui.set("devices", devices);
                } catch (err) {
                    utils.throwAlert({ message: err, variant: "error" });
                }
            }),

            ws.events.on("message", async (data) => {
                // TODO: Need `[]Device` type like data from the websocket
                //this.store.ui.set("devices", devices)
            })
        );

        this.fetchApiDevices().then((devices) => {
            if (!!devices) this.store.ui.set("devices", devices);
        });
    }

    disconnectedCallback() {
        super.disconnectedCallback();
        this.cleanup.run();
    }

    /**
     * @private
     * @returns {Promise<Device[]>}
     */
    async fetchApiDevices() {
        return await api.Get(this.store, "/api/devices");
    }
}

console.debug(`Register the "picow-devices-page"`);
customElements.define("picow-devices-page", PicowDevicesPage);
