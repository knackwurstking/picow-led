import svgOptions from "ui/src/svg/smoothie-line-icons/more-vertical";

import { html, UIIconButton } from "ui";
import { DialogDeviceSetup } from "../../components";
import { utils } from "../../lib";

export class DeviceItemOptions extends UIIconButton {
    static register = () => {
        console.debug(`Register "device-item-options" component`);
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

        this.onclick = async (ev) => {
            ev.stopPropagation();

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
                try {
                    const s = this.uiStore.ui.get("server");
                    const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                    const url = `${
                        s.ssl ? "https:" : "http:"
                    }//${addr}/api/device`;
                    const r = await fetch(url, {
                        method: "DELETE",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(device),
                    });

                    if (!r.ok) {
                        r.text().then((r) => {
                            const message = `Server response to ${url}: ${r}`;
                            utils.throwAlert({ message, variant: "error" });
                            console.error(message);
                        });

                        const message = `Fetch from "${url}" with status code ${r.status}`;
                        console.error(message);
                        utils.throwAlert({ message, variant: "error" });
                    }
                } catch (ex) {
                    utils.throwAlert({ message: ex, variant: "error" });
                }
            });

            dialog.ui.events.on("submit", async (device) => {
                try {
                    const s = this.uiStore.ui.get("server");
                    const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                    const url = `${
                        s.ssl ? "https:" : "http:"
                    }//${addr}/api/device`;
                    const r = await fetch(url, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify(device),
                    });

                    if (!r.ok) {
                        r.text().then((r) => {
                            const message = `Server response to ${url}: ${r}`;
                            utils.throwAlert({ message, variant: "error" });
                            console.error(message);
                        });

                        const message = `Fetch from "${url}" with status code ${r.status}`;
                        console.error(message);
                        utils.throwAlert({ message, variant: "error" });
                    }
                } catch (ex) {
                    utils.throwAlert({ message: ex, variant: "error" });
                }
            });

            document.body.appendChild(dialog);
            dialog.ui.open(true);
        };
    }
}

DeviceItemOptions.register();
