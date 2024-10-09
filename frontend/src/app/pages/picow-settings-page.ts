import { CleanUp, html, UICheck, UIInput, UIStackLayoutPage } from "ui";
import type { PicowStore } from "../../types";

export default class PicowSettingsPage extends UIStackLayoutPage {
    store: PicowStore;
    cleanup: CleanUp;

    constructor() {
        super("settings");

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

        // TODO: Add the new theme picker section select between "original" and gruvbox.
        //       Use the new (ui) Dropdown component
        this.innerHTML = html`
            <ui-flex-grid gap="0.25rem">
                <ui-flex-grid-item>
                    <ui-label primary="Use SSL connections" ripple>
                        <ui-check name="ssl" slot="inputs"></ui-check>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Host">
                        <ui-input
                            name="host"
                            slot="inputs"
                            value="${this.store.ui.get("server").host}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Port">
                        <ui-input
                            name="port"
                            slot="inputs"
                            type="number"
                            value="${this.store.ui.get("server").port}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Theme">
                        <ui-dropdown slot="inputs">
                            <ui-dropdown-options vlaue="original">
                                Original
                            </ui-dropdown-options>

                            <ui-dropdown-options value="gruvbox">
                                Gruvbox
                            </ui-dropdown-options>
                        </ui-dropdown>
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
                const input =
                    this.querySelector<UICheck>(`ui-check[name="ssl"]`);

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
                const input = this.querySelector<UIInput>(
                    `ui-input[name="host"]`
                );

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
                const input = this.querySelector<UIInput>(
                    `ui-input[name="port"]`
                );

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
