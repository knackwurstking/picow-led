import { html, UIIconButton } from "ui";
import { power as svgPower } from "ui/svg/smoothie-line-icons";
import { api } from "../../../lib";

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

        this.innerHTML = svgPower;

        this.ui.events.on("click", async (ev) => {
            ev.stopPropagation();
            if (!this.device) return;

            const prevStateBackup = this.picow.state;
            this.picow.state = "pending";

            try {
                // TODO: Using the device color here, this field
                //       is currently missing
                const color = this.picow.isOn()
                    ? this.device.color.map(() => 0)
                    : [255, 255, 255, 255];

                const ok = await api.Post(this.store, "/api/device/color", {
                    server: {
                        addr: this.device.server.addr,
                    },
                    color,
                });
                if (ok) this.device.color = color;
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
