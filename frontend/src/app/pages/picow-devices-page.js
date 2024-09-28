import { CleanUp } from "ui";
import { html, UIStackLayoutPage } from "ui";

export default class PicowDevicesPage extends UIStackLayoutPage {
    /**
     * @param {object} options
     * @param {PicowStore} options.store
     */
    constructor({ store }) {
        super("devices");

        this.store = store;
        this.cleanup = new CleanUp();

        this.#render();
    }

    #render() {
        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    padding-top: var(--ui-app-bar-height);
                    overflow: auto;
                }
            </style>
        `;

        this.innerHTML = html`
            <ul style="border-radius: var(--ui-radius);"></ul>
        `;
    }

    connectedCallback() {
        super.connectedCallback();

        // -------------------- //
        // Handle AppBar events //
        // -------------------- //
        // TODO: ...
        // ------------------- //
        // Handle Store events //
        // ------------------- //
        // TODO: ...
        // ----------------------- //
        // Handle WebSocket events //
        // ----------------------- //
        // TODO: ...
    }

    disconnectedCallback() {
        super.disconnectedCallback();
        this.cleanup.run();
    }
}

console.debug(`Register the "picow-devices-page"`);
customElements.define("picow-devices-page", PicowDevicesPage);
