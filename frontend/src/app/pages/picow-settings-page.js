import { CleanUp } from "ui";
import { html, UIStackLayoutPage } from "ui";

export default class PicowSettingsPage extends UIStackLayoutPage {
    constructor() {
        super("devices");

        /**
         * @type {PicowStore}
         */
        this.store = document.querySelector(`ui-store`);
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
            <ui-flex-grid gap="0.25rem">
                <ui-flex-grid-item>
                    <ui-label primary="Use SSL connections" ripple>
                        <ui-check name="ssl" slot="input"></ui-check>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Host">
                        <ui-input
                            name="host"
                            slot="input"
                            value="${this.store.ui.get("server").host}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Port">
                        <ui-input
                            name="port"
                            slot="input"
                            type="number"
                            value="${this.store.ui.get("server").port}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>
            </ui-flex-grid>
        `;

        // ------------- //
        // Handle Inputs //
        // ------------- //

        {
            /**
             * @type {NodeJS.Timeout | null}
             */
            let inputEventTimeout = null;
            const resetInputEventTimeout = () => {
                if (!inputEventTimeout) return;
                clearTimeout(inputEventTimeout);
                inputEventTimeout = null;
            };

            // ------------ //
            // SSL Checkbox //
            // ------------ //

            {
                /**
                 * @type {import("ui").UICheck}
                 */
                const input = this.querySelector(`ui-check[name="ssl"]`);
                input.ui.checked = this.store.ui.get("server").ssl;
                input.ui.events.on("input", async (state) => {
                    resetInputEventTimeout();

                    inputEventTimeout = setTimeout(() => {
                        this.store.ui.update("server", (server) => {
                            server.ssl = state;
                            return server;
                        });
                    }, 250);
                });
            }

            // ---------- //
            // Host Input //
            // ---------- //

            {
                /**
                 * @type {import("ui").UIInput}
                 */
                const input = this.querySelector(`ui-check[name="host"]`);
                input.ui.events.on("input", async (host) => {
                    resetInputEventTimeout();

                    inputEventTimeout = setTimeout(() => {
                        this.store.ui.update("server", (server) => {
                            server.host = host;
                            return server;
                        });
                    }, 250);
                });
            }

            // ---------- //
            // Port Input //
            // ---------- //

            {
                /**
                 * @type {import("ui").UIInput}
                 */
                const input = this.querySelector(`ui-check[name="port"]`);
                input.ui.events.on("input", async (port) => {
                    resetInputEventTimeout();

                    inputEventTimeout = setTimeout(() => {
                        this.store.ui.update("server", (server) => {
                            server.port = port;
                            return server;
                        });
                    }, 250);
                });
            }
        }
    }
}

console.debug(`Register the "picow-settings-page"`);
customElements.define("picow-settings-page", PicowSettingsPage);
