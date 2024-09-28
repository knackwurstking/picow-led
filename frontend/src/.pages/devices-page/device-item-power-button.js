import svgPower from "ui/src/svg/smoothie-line-icons/power";

import { html, UIIconButton } from "ui";
import { utils } from "../../lib";

/**
 * HTML: `device-item-power-button`
 *
 * Attribute:
 *  - state: "active" | "pending" [optional]
 */
export class DeviceItemPowerButton extends UIIconButton {
    static register = () => {
        console.debug(`Register "device-item-power-button" component`);
        customElements.define(
            "device-item-power-button",
            DeviceItemPowerButton
        );
    };

    constructor() {
        super();

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        /** @type {Device | null} */
        this.data = null;

        this.picow = {
            root: this,

            get state() {
                return this.root.getAttribute("active");
            },

            set state(state) {
                if (!state) {
                    this.root.removeAttribute("state");
                    return;
                }

                this.root.setAttribute("state", state);
            },

            /** @param {Device | null} d */
            set(d) {
                this.root.data = d;
                this.root.handleButtonColor();
            },

            isOn() {
                return !!this.root.data?.color?.find((n) => n > 0);
            },
        };

        this.render();
    }

    render() {
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

        this.onclick = async (ev) => {
            ev.stopPropagation();

            if (!this.data) return;

            const prevState = this.picow.state;
            this.picow.state = "pending";

            // TODO: Using device color here
            const newColor = this.picow.isOn()
                ? this.data.color.map(() => 0)
                : [255, 255, 255, 255];

            try {
                const s = this.uiStore.ui.get("server");
                const addr = !s.port ? s.host : `${s.host}:${s.port}`;
                const url = `${
                    s.ssl ? "https:" : "http:"
                }//${addr}/api/device/color`;
                const r = await fetch(url, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify({
                        server: {
                            addr: this.data.server.addr,
                        },
                        color: newColor,
                    }),
                });

                if (r.ok) {
                    this.data.color = newColor;
                } else {
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
            } finally {
                this.picow.state = prevState;
            }

            this.handleButtonColor();
        };
    }

    /** @private */
    handleButtonColor() {
        if (this.picow.isOn()) {
            this.picow.state = "active";
        } else {
            this.picow.state = null;
        }
    }
}

DeviceItemPowerButton.register();
