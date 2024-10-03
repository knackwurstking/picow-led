import "../components/status-led";

import { UIButton, UIDrawer, UILabel, html } from "ui";
import type StatusLED from "../components/status-led";
import ws from "../lib/websocket";

export interface Drawer {
    element: UIDrawer;
    open(): void;
    close(): void;
}

export default async function (): Promise<Drawer> {
    const stackLayout = document.querySelector(
        "ui-stack-layout"
    ) as PicowStackLayout;

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

    const devicesButton = el.querySelector(
        `ui-button[name="devices"]`
    ) as UIButton;

    devicesButton.ui.events.on("click", () => {
        stackLayout.ui.clear();
        stackLayout.ui.set("devices");
        el.ui.open = false;
    });

    // --------------- //
    // Settings Button //
    // --------------- //

    const settingsButton = el.querySelector(
        `ui-button[name="settings"]`
    ) as UIButton;

    settingsButton.ui.events.on("click", () => {
        stackLayout.ui.clear();
        stackLayout.ui.set("settings");
        el.ui.open = false;
    });

    // -------------------------------------------- //
    // Setup "/ws" event handler for the status led //
    // -------------------------------------------- //

    {
        const statusLED = el.querySelector(
            `status-led[name="/ws"]`
        ) as StatusLED;

        const label = el.querySelector(`ui-label[name="/ws"]`) as UILabel;
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
