import { customElement } from "lit/decorators.js";

import { css, html, LitElement, type PropertyValues } from "lit";
import {
    Events,
    globalStylesToShadowRoot,
    svg,
    UIButton,
    UIDrawer,
    UIStackLayout,
    UIThemeHandler,
} from "ui";

import * as app from "@app";
import * as lib from "@lib";
import * as types from "@types";

@customElement("picow-app")
class PicowApp extends LitElement {
    private events: Events<types.AppBarEvents> = new Events();

    private queryStore(): types.PicowStore {
        return document.querySelector<types.PicowStore>(`ui-store`)!;
    }

    private queryThemeHandler(): UIThemeHandler {
        return document.querySelector<UIThemeHandler>(`ui-theme-handler`)!;
    }

    private queryAppBar(): types.PGAppBar {
        return this.shadowRoot!.querySelector<types.PGAppBar>(`ui-app-bar`)!;
    }

    private queryDrawer(): UIDrawer {
        return this.shadowRoot!.querySelector<UIDrawer>(`ui-drawer`)!;
    }

    private queryStackLayout(): UIStackLayout<types.PicowStackLayoutPage> {
        return this.shadowRoot!.querySelector<UIStackLayout<types.PicowStackLayoutPage>>(
            `ui-stack-layout`,
        )!;
    }

    static get styles() {
        return css`
            :host {
                position: fixed !important;
                top: 0;
                right: 0;
                bottom: 0;
                left: 0;
            }
        `;
    }

    protected render() {
        return html`
            <div class="is-container no-scrollbar" style="height: 100%;">
                <ui-stack-layout></ui-stack-layout>
            </div>

            ${this.renderAppBar()} ${this.renderDrawer()}
        `;
    }

    protected renderAppBar() {
        return html`
            <ui-app-bar position="top">
                <ui-app-bar-item name="menu" slot="left">
                    <ui-icon-button
                        ghost
                        ripple
                        @click=${(ev: Event) => {
                            this.events.dispatch("menu", ev);
                        }}
                    >
                        ${svg.smoothieLineIcons.menu}
                    </ui-icon-button>
                </ui-app-bar-item>

                <ui-app-bar-item name="status" slot="left">
                    <ui-flex-grid-row align="flex-end" gap="0.25rem">
                        <picow-status-led></picow-status-led>
                        <ui-secondary style="white-space: nowrap;"> Offline </ui-secondary>
                    </ui-flex-grid-row>
                </ui-app-bar-item>

                <ui-app-bar-item name="title" slot="center">
                    <h4 style="white-space: nowrap;">PicoW LED</h4>
                </ui-app-bar-item>

                <ui-app-bar-item name="add" slot="right" hidden>
                    <ui-icon-button
                        ghost
                        ripple
                        @click=${(ev: Event) => {
                            this.events.dispatch("add", ev);
                        }}
                    >
                        ${svg.smoothieLineIcons.plus}
                    </ui-icon-button>
                </ui-app-bar-item>
            </ui-app-bar>
        `;
    }

    protected renderDrawer() {
        return html`
            <ui-drawer>
                <ui-drawer-group no-fold>
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
            </ui-drawer>
        `;
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        globalStylesToShadowRoot(this.shadowRoot!);

        this.initializeStore();
        this.initializeStackLayout();
        this.initializeAppBar();
        this.initializeDrawer();

        setTimeout(async () => {
            const store = this.queryStore();
            const stackLayout = this.queryStackLayout();
            const drawer = this.queryDrawer();

            // Set the start page
            if (!!store.getData("currentPage")) {
                const currentPage = store.getData("currentPage");
                if (currentPage !== undefined && currentPage !== null) {
                    stackLayout.setPage(currentPage);
                }
            } else {
                drawer.open = true;
            }

            // Handle websocket error messages
            lib.ws.events.addListener("message-error", (msg) =>
                lib.throwAlert({ message: msg, variant: "error" }),
            );

            // Updating the websocket server if store data "server" changes
            store.addListener(
                "server",
                async (server) => {
                    lib.ws.server = server;
                },
                true,
            );
        });
    }

    private initializeStore() {
        const store = this.queryStore();
        const themeHandler = this.queryThemeHandler();

        themeHandler.theme = store.getData("currentTheme")?.theme || themeHandler.theme;

        store.setData("currentPage", "", true);
        store.setData("devices", [], true);
        store.setData("devicesColor", {}, true);

        store.setData(
            "server",
            {
                ssl: !!location.protocol.match(/(https)/),
                host: location.hostname,
                port: location.port,
            },
            true,
        );
    }

    private initializeStackLayout() {
        const store = this.queryStore();
        const appBar = this.queryAppBar();
        const drawer = this.queryDrawer();

        const stackLayout =
            this.shadowRoot!.querySelector<UIStackLayout<types.PicowStackLayoutPage>>(
                `ui-stack-layout`,
            )!;

        stackLayout.registerPage("devices", async () => {
            const page = new app.PicowDevicesPage();
            page.picowAppEvents = this.events;
            return page;
        });

        stackLayout.registerPage("settings", async () => {
            return new app.PicowSettingsPage();
        });

        stackLayout.events.addListener("change", async ({ current }) => {
            const addItem = appBar.contentName("add")!;

            // Reset all layouts (AppBar buttons and title)
            appBar.title = "PicoW LED";
            addItem.hide();

            if (!current) {
                drawer.open = true;
                return;
            }

            store.setData("currentPage", current.name as types.PicowStackLayoutPage);
            switch (current.name) {
                case "devices":
                    store.setData("currentPage", current.name);
                    appBar.title = "Devices";
                    addItem.show();
                    break;

                case "settings":
                    store.setData("currentPage", current.name);
                    appBar.title = "Settings";
                    break;

                default:
                    appBar.title = current.name;
                    break;
            }
        });
    }

    private initializeAppBar() {
        const appBar = this.queryAppBar();
        const drawer = this.queryDrawer();

        this.events.addListener("menu", () => {
            drawer.open = true;
        });

        const statusItem = appBar.contentName("status")!;
        const picowStatusLED = statusItem.querySelector<app.PicowStatusLED>(`picow-status-led`)!;

        const text = statusItem.querySelector(`ui-secondary`)!;

        lib.ws.events.addListener("open", () => {
            picowStatusLED.active = true;
            text.innerHTML = "Online";
        });

        lib.ws.events.addListener("close", () => {
            picowStatusLED.active = false;
            text.innerHTML = "Offline";
        });
    }

    private initializeDrawer() {
        const drawer = this.queryDrawer();
        const stackLayout = this.queryStackLayout();

        let button: UIButton;
        const pages: types.PicowStackLayoutPage[] = ["devices", "settings"];
        for (const name of pages) {
            button = this.shadowRoot!.querySelector(`ui-button[name="${name}"]`)!;
            button.onclick = () => {
                stackLayout.clearStack();
                stackLayout.setPage(name);
                drawer.open = false;
            };
        }
    }
}

export default PicowApp;
