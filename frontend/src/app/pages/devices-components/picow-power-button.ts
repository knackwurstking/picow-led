import { css, html, LitElement } from "lit";
import { customElement, property } from "lit/decorators.js";
import { svg } from "ui";
import { ws, WSEventsDevice } from "../../../lib/websocket";
import { PicowStore } from "../../../types";

/**
 * **Tag**: picow-power-button
 *
 * **Attributes**:
 *  - device: `WSEventsDevice` - [json]
 *  - state: `"active" | "pending" | ""`
 */
@customElement("picow-power-button")
export class PicowPowerButton extends LitElement {
    @property({ type: Object, attribute: "device", reflect: true })
    device?: WSEventsDevice;

    @property({ type: String, attribute: "state", reflect: true })
    state: "active" | "pending" | "" = "";

    store: PicowStore = document.querySelector(`ui-store`)!;

    static get styles() {
        return css`
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
        `;
    }

    protected render() {
        return html`
            <ui-icon-button
                ghost
                @click=${async (ev: MouseEvent) => {
                    ev.stopPropagation();
                    if (!this.device) return;

                    const prevStateBackup = this.state;
                    this.state = "pending";

                    try {
                        const color: number[] = this.isOn()
                            ? this.device.color!.map(() => 0) // Turn OFF
                            : this.store.getData("devicesColor")?.[
                                  this.device.server.addr
                              ] || [255, 255, 255, 255]; // Turn ON

                        await ws.request("POST api.device.color", {
                            addr: this.device.server.addr,
                            color: color,
                        });
                    } finally {
                        this.state = prevStateBackup;
                    }

                    this.updateColor();
                }}
            >
                ${svg.smoothieLineIcons.power}
            </ui-icon-button>
        `;
    }

    attributeChangedCallback(
        name: string,
        _old: string | null,
        value: string | null
    ): void {
        super.attributeChangedCallback(name, _old, value);

        switch (name) {
            case "device":
                this.updateColor();
                break;
        }
    }

    private isOn() {
        return !!this.device?.color?.find((n: number) => n > 0);
    }

    private updateColor() {
        if (this.isOn()) this.state = "active";
        else this.state = "";
    }
}
