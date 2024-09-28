import { CleanUp, html, styles } from "ui";
import { deviceEvents } from "../../lib";
import { DeviceItemOfflineMarker } from "./device-item-offline-marker";
import { DeviceItemOptions } from "./device-item-options";
import { DeviceItemPowerButton } from "./device-item-power-button";

export class DeviceItem extends HTMLElement {
    static register = () => {
        console.debug(`Register "device-item" component`);
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

                if (!!device.color) {
                    this.root.style.setProperty(
                        "--current-color",
                        `rgb(${device.color[0] || 0}, ${
                            device.color[1] || 0
                        }, ${device.color[2] || 0})`
                    );
                }
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
                    display: block !important;
                    position: relative;
                    border-radius: var(--ui-radius) !important;
                }

                div.current-color {
                    position: absolute;
                    top: var(--ui-spacing);
                    right: var(--ui-spacing);
                    bottom: var(--ui-spacing);
                    left: var(--ui-spacing);
                    border-radius: var(--ui-radius);
                    box-shadow: 0 0 8px 1px var(--current-color, transparent);
                    transition: box-shadow 0.35s linear;
                }
            </style>

            <div class="current-color"></div>
            <slot></slot>
        `;

        this.innerHTML = html`
            <li class="is-card" style="cursor: pointer;">
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

        /**
         * @type {HTMLElement}
         */
        const card = this.querySelector("li.is-card");
        card.onclick = async () => {
            // TODO: Open a color picker dialog to select a color
        };

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

DeviceItem.register();
