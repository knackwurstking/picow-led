import "../components/status-led";

import { UIDrawer, html } from "ui";
import ws from "../lib/websocket";

/**
 * @typedef {{
 *  element: UIDrawer;
 *  open(): void;
 *  close(): void;
 * }} Drawer
 */

/**
 * @param {object} options
 * @param {import("ui").UIStackLayout<PicowStackLayout_Pages>} options.stackLayout
 * @returns {Promise<Drawer>}
 */
export default async function ({ stackLayout }) {
    const el = new UIDrawer();

    el.innerHTML = html`
        <ui-drawer-group nofold>
            <ui-drawer-group-item style="column">
                <ui-label name="/ws" style="width: 100%;">
                    <status-led name="/ws"></status-led>
                </ui-label>
            </ui-drawer-group-item>
        </ui-drawer-group>

        <ui-drawer-group nofold>
            <ui-drawer-group-item>
                <ui-button
                    name="devices"
                    style="justify-content: flex-start;"
                    variant="ghost"
                    color="primary"
                >
                    Devices
                </ui-button>
            </ui-drawer-group-item>

            <ui-drawer-group-item>
                <ui-button
                    name="groups"
                    style="justify-content: flex-start;"
                    variant="ghost"
                    color="primary"
                    disabled
                >
                    Groups
                </ui-button>
            </ui-drawer-group-item>

            <ui-drawer-group-item>
                <ui-button
                    name="scenes"
                    style="justify-content: flex-start;"
                    variant="ghost"
                    color="primary"
                    disabled
                >
                    Scenes
                </ui-button>
            </ui-drawer-group-item>

            <ui-drawer-group-item>
                <ui-button
                    name="settings"
                    style="justify-content: flex-start;"
                    variant="ghost"
                    color="primary"
                >
                    Settings
                </ui-button>
            </ui-drawer-group-item>
        </ui-drawer-group>
    `;

    // -------------- //
    // Devices Button //
    // -------------- //

    /**
     * @type {import("ui").UIButton}
     */
    const devicesButton = el.querySelector(`ui-button[name="devices"]`);
    devicesButton.ui.events.on("click", () => {
        stackLayout.ui.clear();
        stackLayout.ui.set("devices");
        el.ui.open = false;
    });

    // --------------- //
    // Settings Button //
    // --------------- //

    /**
     * @type {import("ui").UIButton}
     */
    const settingsButton = el.querySelector(`ui-button[name="settings"]`);
    settingsButton.ui.events.on("click", () => {
        stackLayout.ui.clear();
        stackLayout.ui.set("settings");
        el.ui.open = false;
    });

    // -------------------------------------------- //
    // Setup "/ws" event handler for the status led //
    // -------------------------------------------- //

    {
        /**
         * @type {import("../components").StatusLED}
         */
        const statusLED = el.querySelector(`status-led[name="/ws"]`);

        /**
         * @type {import("ui").UILabel}
         */
        const label = el.querySelector(`ui-label[name="/ws"]`);
        label.ui.primary = ws.path;
        label.ui.secondary = ws.origin;

        ws.events.on("server", async () => {
            label.ui.primary = ws.path;
            label.ui.secondary = ws.origin;
        });
        ws.events.on("close", () => statusLED.removeAttribute("active"));
        ws.events.on("open", () => statusLED.setAttribute("active", ""));
    }

    return {
        element: el,

        open() {
            el.ui.open = true;
        },

        close() {
            el.ui.open = false;
        },
    };
}
