import { html, UIDrawer } from "ui";
import { deviceEvents, devicesEvents } from "../../lib";
import { StatusLED } from "../status-led";

export class PicowDrawer extends UIDrawer {
    static register = () => {
        StatusLED.register();

        customElements.define("picow-drawer", PicowDrawer);
    };

    constructor() {
        super();

        /** @type {PicowStore} */
        this.uiStore = document.querySelector("ui-store");

        /** @type {PicowStackLayout} */
        this.uiLayout = document.querySelector("ui-stack-layout");

        this.render();
    }

    render() {
        this.innerHTML = html`
            <ui-drawer-group nofold>
                <ui-drawer-group-item style="column">
                    <ui-label name="device" style="width: 100%;">
                        <status-led name="device"></status-led>
                    </ui-label>
                </ui-drawer-group-item>

                <ui-drawer-group-item style="column">
                    <ui-label name="devices" style="width: 100%;">
                        <status-led name="devices"></status-led>
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

        this.renderButtons();
        this.renderStatusLEDs();
    }

    /** @private */
    renderButtons() {
        /**
         * Devices Page Button
         *
         * @type {import("ui").UIButton}
         */
        const devicesButton = this.querySelector(`ui-button[name="devices"]`);
        devicesButton.ui.events.on("click", () => {
            this.uiLayout.ui.clear();
            this.uiLayout.ui.set("devices");
            this.ui.open = false;
        });

        /**
         * Settings Page Button
         *
         * @type {import("ui").UIButton}
         */
        const settingsButton = this.querySelector(`ui-button[name="settings"]`);
        settingsButton.ui.events.on("click", () => {
            this.uiLayout.ui.clear();
            this.uiLayout.ui.set("settings");
            this.ui.open = false;
        });
    }

    /** @private */
    renderStatusLEDs() {
        {
            /**
             * "/events/device"
             *
             * @type {import("../status-led").StatusLED}
             */
            const deviceStatusLED = this.querySelector(
                `status-led[name="device"]`
            );

            deviceEvents.events.on("server", async () => {
                /**
                 * @type {import("ui").UILabel}
                 */
                const label = this.querySelector(`ui-label[name="device"]`);
                label.ui.primary = deviceEvents.path;
                label.ui.secondary = deviceEvents.origin;
            });

            deviceEvents.events.on("close", () => {
                deviceStatusLED.removeAttribute("active");
            });

            deviceEvents.events.on("open", () => {
                deviceStatusLED.setAttribute("active", "");
            });
        }

        {
            /**
             * "/events/devices"
             *
             * @type {import("../status-led").StatusLED}
             */
            const devicesStatusLED = this.querySelector(
                `status-led[name="devices"]`
            );

            devicesEvents.events.on("server", async () => {
                /**
                 * @type {import("ui").UILabel}
                 */
                const label = this.querySelector(`ui-label[name="devices"]`);
                label.ui.primary = devicesEvents.path;
                label.ui.secondary = devicesEvents.origin;
            });

            devicesEvents.events.on("close", () => {
                devicesStatusLED.removeAttribute("active");
            });

            devicesEvents.events.on("open", () => {
                devicesStatusLED.setAttribute("active", "");
            });
        }
    }
}
