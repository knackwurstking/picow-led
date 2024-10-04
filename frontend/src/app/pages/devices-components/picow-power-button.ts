import { power as svgPower } from "ui/svg/smoothie-line-icons";

import { html, UIIconButton } from "ui";
import * as api from "../../../lib/api";
import type { WSEvents_Device } from "../../../lib/websocket/ws-events";
import type { PicowStore } from "../../../types";

export type PicowPowerButton_States = "active" | "pending" | null;

class PicowPowerButton_Picow {
    root: PicowPowerButton;

    constructor(root: PicowPowerButton) {
        this.root = root;
    }

    get state(): PicowPowerButton_States {
        return this.root.getAttribute("active") as PicowPowerButton_States;
    }

    set state(state: PicowPowerButton_States) {
        if (!state) {
            this.root.removeAttribute("state");
            return;
        }

        this.root.setAttribute("state", state);
    }

    set(device: WSEvents_Device) {
        this.root.device = device;
        this.root.updateColor();
    }

    isOn() {
        return !!this.root.device?.color?.find((n: number) => n > 0);
    }
}

export default class PicowPowerButton extends UIIconButton {
    store: PicowStore;
    device: WSEvents_Device | null;
    picow: PicowPowerButton_Picow;

    constructor() {
        super();

        this.store = document.querySelector(`ui-store`);
        this.device = null;

        this.picow = new PicowPowerButton_Picow(this);

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
                // TODO: Using the (stored) device color here
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

            this.updateColor();
        });
    }

    updateColor() {
        if (this.picow.isOn()) this.picow.state = "active";
        else this.picow.state = null;
    }
}

console.debug(`Register the "picow-power-button"`);
customElements.define("picow-power-button", PicowPowerButton);
