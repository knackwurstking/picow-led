import { power as svgPower } from "ui/svg/smoothie-line-icons";
import { html, UIIconButton } from "ui";
import { api, utils } from "../../../lib";

export default class PicowPowerButton extends UIIconButton {
    constructor() {
        super();

        /**
         * @type {PicowStore}
         */
        this.store = document.querySelector(`ui-store`);

        /**
         * @type {Device | null}
         */
        this.device = null;

        this.picow = {
            root: this,

            get state() {
                // @ts-expect-error
                return this.root.getAttribute("active");
            },

            /**
             * @param {"pending" | "active" | null} state
             */
            set state(state) {
                if (!state) {
                    this.root.removeAttribute("state");
                    return;
                }

                this.root.setAttribute("state", state);
            },

            /**
             * @param {Device} device
             */
            set(device) {
                this.root.device = device;
                this.root.#handleButtonColor();
            },

            isOn() {
                return !!this.root.device?.color?.find((n) => n > 0);
            },
        };

        this.#render();
    }

    #render() {
        this.ui.noripple = true;
        this.ui.ghost = true;

        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    height: 100%;
                    width: 3rem;
                    color: black;
                }

                :host([state="active"]) {
                    color: rgb(0, 255, 0);
                }

                :host([state="pending"]) {
                    color: yellow;
                }
            </style>
        `;

        this.inneHTML = svgPower;

        this.ui.events.on("click", async (ev) => {
            ev.stopPropagation();
            if (!this.device) return;

            const prevStateBackup = this.picow.state;
            this.picow.state = "pending";

            const url = await api.url(this.store, "/api/device/color");

            try {
                // TODO: Using the device color here, this field is currently missing
                const color = this.picow.isOn()
                    ? this.device.color.map(() => 0)
                    : [255, 255, 255, 255];

                const resp = await fetch(url, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        server: {
                            addr: this.device.server.addr,
                        },
                        color,
                    }),
                });

                if (resp.ok) {
                    this.device.color = color;
                } else {
                    resp.text().then((r) => {
                        const m = `Server response to ${url}: ${r}`;
                        utils.throwAlert({ message: m, variant: "error" });
                        console.error(m);
                    });

                    const m = `Fetch from ${url} with status code ${resp.status}`;
                    console.error(m);
                    utils.throwAlert({ message: m, variant: "error" });
                }
            } catch (err) {
                utils.throwAlert({ message: err, variant: "error" });
                console.error(err);
            } finally {
                this.picow.state = prevStateBackup;
            }

            this.#handleButtonColor();
        });
    }

    #handleButtonColor() {
        if (this.picow.isOn()) this.picow.state = "active";
        else this.picow.state = null;
    }
}

console.debug(`Register the "picow-power-button"`);
customElements.define("picow-power-button", PicowPowerButton);
