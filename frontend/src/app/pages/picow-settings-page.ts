import { CleanUp, html, UICheck, UIInput, UIStackLayoutPage } from "ui";

export default class PicowSettingsPage extends UIStackLayoutPage {
    store: PicowStore;
    cleanup: CleanUp;

    constructor() {
        super("devices");

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
            let inputEventTimeout: NodeJS.Timeout | null = null;
            const resetInputEventTimeout = () => {
                if (!inputEventTimeout) return;
                clearTimeout(inputEventTimeout);
                inputEventTimeout = null;
            };

            // ------------ //
            // SSL Checkbox //
            // ------------ //

            {
                const input = this.querySelector(
                    `ui-check[name="ssl"]`
                ) as UICheck;

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
                const input = this.querySelector(
                    `ui-input[name="host"]`
                ) as UIInput;

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
                const input = this.querySelector(
                    `ui-input[name="port"]`
                ) as UIInput;

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
