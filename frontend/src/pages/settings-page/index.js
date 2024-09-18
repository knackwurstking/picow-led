import { html, UIStackLayoutPage } from "ui";

export class SettingsPage extends UIStackLayoutPage {
    static register = () => {
        customElements.define("settings-page", SettingsPage);
    };

    constructor() {
        super("settings");

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        this.render();
    }

    shadowRender() {
        super.shadowRender();
        this.classList.add("no-scrollbar");
        this.shadowRoot.innerHTML += html`
            <style>
                :host {
                    padding-top: var(--ui-app-bar-height);
                    overflow: auto;
                }
            </style>
        `;
    }

    render() {
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
                            value="${this.uiStore.ui.get("server").host}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Port">
                        <ui-input
                            name="port"
                            slot="input"
                            value="${this.uiStore.ui.get("server").port}"
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>
            </ui-flex-grid>
        `;

        this.renderInputs();
    }

    renderInputs() {
        /** @type {NodeJS.Timeout | null} */
        let timeout = null;

        const handleTimeout = () => {
            if (!!timeout) {
                clearTimeout(timeout);
                timeout = null;
            }
        };

        /** @type {import("ui").UICheck} */
        const ssl = this.querySelector(`ui-check[name="ssl"]`);
        ssl.ui.events.on("input", (state) => {
            handleTimeout();

            timeout = setTimeout(() => {
                this.uiStore.ui.update("server", (server) => {
                    return {
                        ...server,
                        ssl: state,
                    };
                });
            }, 250);
        });

        /** @type {import("ui").UIInput} */
        const hostInput = this.querySelector(`ui-input[name="host"]`);
        hostInput.ui.events.on("input", (value) => {
            handleTimeout();

            timeout = setTimeout(() => {
                this.uiStore.ui.update("server", (server) => {
                    return {
                        ...server,
                        host: value,
                    };
                });
            }, 250);
        });

        /** @type {import("ui").UIInput} */
        const portInput = this.querySelector(`ui-input[name="port"]`);
        portInput.ui.events.on("input", (value) => {
            handleTimeout();

            timeout = setTimeout(() => {
                this.uiStore.ui.update("server", (server) => {
                    return {
                        ...server,
                        port: value,
                    };
                });
            }, 250);
        });
    }
}
