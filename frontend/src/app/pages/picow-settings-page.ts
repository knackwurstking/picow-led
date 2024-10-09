import {
    CleanUp,
    html,
    UICheck,
    UIInput,
    UISelect,
    UIStackLayoutPage,
    UIThemeHandler,
} from "ui";
import type { PicowStore } from "../../types";
import type { UIThemeHandler_Theme } from "ui/src/ui-theme-handler/ui-theme-handler";

export default class PicowSettingsPage extends UIStackLayoutPage {
    store: PicowStore;
    cleanup: CleanUp;
    themeHandler: UIThemeHandler;

    constructor() {
        super("settings");

        this.cleanup = new CleanUp();

        this.store = document.querySelector(`ui-store`);
        this.themeHandler = document.querySelector(`ui-theme-handler`);

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
                        <ui-select name="theme" keep-open>
                            <ui-select-option value="original">
                                Original
                            </ui-select-option>

                            <ui-select-option value="gruvbox">
                                Gruvbox
                            </ui-select-option>
                        </ui-select>
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

            // ------------ //
            // Theme Select //
            // ------------ //

            {
                const theme = this.querySelector<UISelect>(
                    `ui-select[name="theme"]`
                );

                theme.ui.options().forEach((option) => {
                    const currentTheme = this.themeHandler.ui.theme;
                    option.ui.selected = option.ui.value === currentTheme;
                });

                theme.ui.events.on("change", async (option) => {
                    this.themeHandler.ui.theme = option.ui
                        .value as UIThemeHandler_Theme;

                    this.store.ui.set("currentTheme", {
                        theme: this.themeHandler.ui.theme,
                    });
                });
            }
        }
    }
}

console.debug(`Register the "picow-settings-page"`);
customElements.define("picow-settings-page", PicowSettingsPage);
