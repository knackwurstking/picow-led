import { html, UIStackLayoutPage } from "ui";

/**
 * @typedef {{
 *  element: UIStackLayoutPage;
 *  get name(): string;
 * }} SettingsPage
 */

/**
 * @param {object} options
 * @param {PicowStore} options.store
 * @returns {Promise<SettingsPage>}
 */
export default async function ({ store }) {
    const pageName = "settings";
    const page = new UIStackLayoutPage(pageName);

    page.style.paddingTop = "var(--ui-app-bar-height)";
    page.style.overflow = "auto";

    page.innerHTML = html`
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
                        value="${store.ui.get("server").host}"
                    ></ui-input>
                </ui-label>
            </ui-flex-grid-item>

            <ui-flex-grid-item>
                <ui-label primary="Server Port">
                    <ui-input
                        name="port"
                        slot="input"
                        type="number"
                        value="${store.ui.get("server").port}"
                    ></ui-input>
                </ui-label>
            </ui-flex-grid-item>
        </ui-flex-grid>
    `;

    /**
     * @type {NodeJS.Timeout | null}
     */
    let inputEventTimeout = null;
    const resetInputEventTimeout = () => {
        if (!inputEventTimeout) return;
        clearTimeout(inputEventTimeout);
        inputEventTimeout = null;
    };

    //
    // SSL Checkbox
    //

    {
        /**
         * @type {import("ui").UICheck}
         */
        const input = page.querySelector(`ui-check[name="ssl"]`);
        input.ui.checked = store.ui.get("server").ssl;
        input.ui.events.on("input", async (state) => {
            resetInputEventTimeout();

            inputEventTimeout = setTimeout(() => {
                store.ui.update("server", (server) => {
                    server.ssl = state;
                    return server;
                });
            }, 250);
        });
    }

    //
    // Host Input
    //

    {
        /**
         * @type {import("ui").UIInput}
         */
        const input = page.querySelector(`ui-check[name="host"]`);
        input.ui.events.on("input", async (host) => {
            resetInputEventTimeout();

            inputEventTimeout = setTimeout(() => {
                store.ui.update("server", (server) => {
                    server.host = host;
                    return server;
                });
            }, 250);
        });
    }

    //
    // Port Input
    //

    {
        /**
         * @type {import("ui").UIInput}
         */
        const input = page.querySelector(`ui-check[name="port"]`);
        input.ui.events.on("input", async (value) => {
            const port = parseInt(value, 10);
            resetInputEventTimeout();

            inputEventTimeout = setTimeout(() => {
                store.ui.update("server", (server) => {
                    server.port = port;
                    return server;
                });
            }, 250);
        });
    }

    return {
        element: page,

        get name() {
            return pageName;
        },
    };
}
