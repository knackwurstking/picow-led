import { CleanUp, globalStylesToShadowRoot, html } from "ui";
import { deviceEvents } from "../../../lib";

export default class PicowDeviceItem extends HTMLElement {
    /**
     * @param {object} options
     * @param {Device | null} [options.device]
     */
    constructor({ device = null }) {
        super();

        /**
         * @type {Device | null}
         */
        this.device = device;
        this.cleanup = new CleanUp();

        this.picow = {
            root: this,

            /**
             * @param {Device} device
             */
            set(device) {
                // TODO: ...
            },
        };

        this.#render();
    }

    #render() {
        this.attachShadow({ mode: "open" });
        globalStylesToShadowRoot(this.shadowRoot);

        // TODO: Continue here...
        this.shadowRoot.innerHTML = html``;

        this.picow.set(this.device);
    }

    connectedCallback() {
        this.cleanup.add(
            deviceEvents.events.on("message", (device) => {
                if (device.server.addr !== this.device.server.addr) return;
                this.picow.set(device);
            })
        );
    }

    disconnectedCallback() {
        this.cleanup.run();
    }
}

customElements.define("picow-device-item", PicowDeviceItem);
