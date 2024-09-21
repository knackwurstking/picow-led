import { CleanUp, html } from "ui";
import { deviceEvents } from "../../lib";
import { DeviceItemOfflineMarker } from "./device-item-offline-marker";
import { DeviceItemOptions } from "./device-item-options";
import { DeviceItemPowerButton } from "./device-item-power-button";

export class DeviceItem extends HTMLElement {
    static register = () => {
        DeviceItemOfflineMarker.register();
        DeviceItemOptions.register();
        DeviceItemPowerButton.register();

        customElements.define("device-item", DeviceItem);
    };

    /**
     * @param {Device | null} [data]
     */
    constructor(data = null) {
        super();

        this.cleanup = new CleanUp();

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        /** @type {Device | null} */
        this.data = data;

        this.picow = {
            root: this,

            /**
             * @param {Device} device
             */
            async set(device) {
                this.root.data = device;

                const li = this.root.querySelector("li");
                li.setAttribute("data-addr", device.server.addr);

                this.root.setPowerButton();
                this.root.setOptionsButton();
                this.root.setOfflineMarker();
                this.root.setText();
            },
        };

        this.render();
    }

    render() {
        this.attachShadow({ mode: "open" });
        this.classList.add("no-user-select");
        this.shadowRoot.innerHTML = html`
            <style>
                :host {
                    display: block;
                }
            </style>

            <slot></slot>
        `;

        this.innerHTML = html`
            <li class="is-card">
                <device-item-offline-marker></device-item-offline-marker>
                <ui-label>
                    <ui-flex-grid-row gap="0.25rem">
                        <ui-flex-grid-item>
                            <device-item-power-button>
                            </device-item-power-button>
                        </ui-flex-grid-item>

                        <ui-flex-grid-item>
                            <device-item-options></device-item-options>
                        </ui-flex-grid-item>
                    </ui-flex-grid-row>
                </ui-label>
            </li>
        `;

        this.picow.set(this.data);
    }

    connectedCallback() {
        this.cleanup.add(
            deviceEvents.events.on("message", (d) => {
                if (d.server.addr !== this.data.server.addr) return;
                this.picow.set(d);
            })
        );
    }

    disconnectedCallback() {
        this.cleanup.run();
    }

    /** @private */
    setPowerButton() {
        /** @type {import("./device-item-power-button").DeviceItemPowerButton} */
        const power = this.querySelector(`device-item-power-button`);
        power.picow.set(this.data);
    }

    /** @private */
    setOptionsButton() {
        /** @type {import("./device-item-options").DeviceItemOptions} */
        const options = this.querySelector(`device-item-options`);
        options.set(this.data);
    }

    /** @private */
    setOfflineMarker() {
        /** @type {import("./device-item-offline-marker").DeviceItemOfflineMarker} */
        const offlineMarker = this.querySelector("device-item-offline-marker");
        offlineMarker.picow.hide = !this.data.server.isOffline;
    }

    /** @private */
    setText() {
        /** @type {import("ui").UILabel} */
        const uiLabel = this.querySelector("ui-label");

        let primary = this.data.server.name || "";
        let secondary = this.data.server.addr;

        if (!primary) {
            primary = this.data.server.addr;
            secondary = "&nbsp;";
        }

        uiLabel.ui.primary = primary;
        uiLabel.ui.secondary = secondary;
    }
}
