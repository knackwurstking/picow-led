import svgOptions from "ui/src/svg/smoothie-line-icons/more-vertical";

import { html, UIIconButton } from "ui";
import { DialogDeviceSetup } from "../../components";

export class DeviceItemOptions extends UIIconButton {
    static register = () => {
        customElements.define("device-item-options", DeviceItemOptions);
    };

    /**
     * @param {Device | null} data
     */
    constructor(data = null) {
        super();

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        /** @type {Device | null} */
        this.data;
        this.set(data);

        this.render();
    }

    /** @param {Device | null} device */
    set(device) {
        this.data = device;
    }

    render() {
        this.ui.ghost = true;

        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    height: 100%;
                }
            </style>
        `;

        this.innerHTML = svgOptions;

        this.ui.events.on("click", () => {
            const dialog = new DialogDeviceSetup({
                name: this.data.server.name,
                addr: this.data.server.addr,
                pins: this.data.pins,
                allowDeletion: true,
            });

            dialog.ui.events.on("close", () => {
                document.body.removeChild(dialog);
            });

            dialog.ui.events.on("delete", async (device) => {
                const s = this.uiStore.ui.get("server");
                const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                const url = `${s.ssl ? "https:" : "http:"}//${addr}/api/device`;
                const r = await fetch(url, {
                    method: "DELETE",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(device),
                });

                if (!r.ok) {
                    // TODO: Add an "error" alert
                    r.text().then((r) => console.error(r));
                    console.error(
                        `Fetch from "${url}" with status code ${r.status}`
                    );
                }
            });

            dialog.ui.events.on("submit", async (device) => {
                const s = this.uiStore.ui.get("server");
                const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                const url = `${s.ssl ? "https:" : "http:"}//${addr}/api/device`;
                const r = await fetch(url, {
                    method: "PUT",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(device),
                });

                if (!r.ok) {
                    // TODO: Add an "error" alert
                    r.text().then((r) => console.error(r));
                    console.error(
                        `Fetch from "${url}" with status code ${r.status}`
                    );
                }
            });

            document.body.appendChild(dialog);
            dialog.ui.open(true);
        });
    }
}
