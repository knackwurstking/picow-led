import { CleanUp, html, UIStackLayoutPage } from "ui";
import { devicesEvents, utils } from "../../lib";

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

            this.appBar.events.on("add", () => {
                // TODO: Create device setup dialog
            }),

            // ------------------- //
            // Handle Store events //
            // ------------------- //

            this.store.ui.on("devices", (devices) => {
                const list = this.querySelector("ul");
                while (!!list.firstChild) list.removeChild(list.firstChild);
                for (const device of devices) {
                    // TODO: Create and append the device item
                }
            }),

            // ----------------------- //
            // Handle WebSocket events //
            // ----------------------- //

            devicesEvents.events.on("open", async () => {
                try {
                    this.store.ui.set("devices", await this.fetchApiDevices());
                } catch (err) {
                    utils.throwAlert({ message: err, variant: "error" });
                }
            }),

            devicesEvents.events.on("message", async (devices) =>
                this.store.ui.set("devices", devices)
            )
        );

        this.fetchApiDevices()
            .then((devices) => {
                this.store.ui.set("devices", devices);
            })
            .catch((err) =>
                utils.throwAlert({ message: err, variant: "error" })
            );
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
        const server = this.store.ui.get("server");
        const addr = !server.port
            ? server.host
            : `${server.host}:${server.port}`;
        const url = `${server.ssl ? "https:" : "http:"}//${addr}/api/devices`;

        const resp = await fetch(url, { method: "GET" });
        if (!resp.ok) {
            resp.text().then((r) => {
                const m = `Server response to ${url}: ${r}`;
                utils.throwAlert({ message: m, variant: "error" });
                console.error(m);
            });

            const m = `Fetch from "${url}" with status code ${resp.status}`;
            utils.throwAlert({ message: m, variant: "error" });
            console.error(m);
            return;
        }

        return await resp.json();
    }
}

console.debug(`Register the "picow-devices-page"`);
customElements.define("picow-devices-page", PicowDevicesPage);
