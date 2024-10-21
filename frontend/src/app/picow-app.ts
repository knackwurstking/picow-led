import { css, html, LitElement, type PropertyValues } from "lit";
import { customElement } from "lit/decorators.js";
import {
    Events,
    globalStylesToShadowRoot,
    svg,
    UIAppBar,
    UIStackLayout,
    UIThemeHandler,
} from "ui";
import { PicowStatusLED } from "../components/picow-status-led";
import { throwAlert } from "../lib/utils";
import { ws } from "../lib/websocket";
import { AppBarEvents, PicowStackLayoutPages, PicowStore } from "../types";
import { PicowDevicesPage } from "./pages/picow-devices-page";
import { PicowSettingsPage } from "./pages/picow-settings-page";
import { PicowDrawer } from "./picow-drawer";

/**
 * **Tag**: picow-app
 */
@customElement("picow-app")
export class PicowApp extends LitElement {
    private events: Events<AppBarEvents> = new Events();

    private drawer: PicowDrawer = new PicowDrawer();

    private store(): PicowStore {
        return document.querySelector<PicowStore>(`ui-store`)!;
    }

    private themeHandler(): UIThemeHandler {
        return document.querySelector<UIThemeHandler>(`ui-theme-handler`)!;
    }

    private appBar(): UIAppBar {
        return this.shadowRoot!.querySelector<UIAppBar>(`ui-app-bar`)!;
    }

    private stackLayout(): UIStackLayout<PicowStackLayoutPages> {
        return this.shadowRoot!.querySelector<
            UIStackLayout<PicowStackLayoutPages>
        >(`ui-stack-layout`)!;
    }

    static get styles() {
        return css`
            :host {
                position: fixed !important;
                top: 0;
                right: 0;
                bottom: 0;
                left: 0;
                overflow: auto;
            }
        `;
    }

    protected render() {
        return html`
            <ui-container style="width: 100%; height: 100%;">
                <ui-stack-layout></ui-stack-layout>
            </ui-container>

            ${this.renderAppBar()} ${this.drawer}
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
                        <ui-secondary style="white-space: nowrap;">
                            Offline
                        </ui-secondary>
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
        return html``; // TODO: ...
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        globalStylesToShadowRoot(this.shadowRoot!);

        this.initializeStore();
        this.initializeStackLayout();
        this.initializeAppBar();
        this.initializeDrawer();

        setTimeout(async () => {
            const store = this.store();
            const stackLayout = this.stackLayout();

            // Set the start page
            if (!!store.getData("currentPage")) {
                const currentPage = store.getData("currentPage");
                if (currentPage !== undefined && currentPage !== null) {
                    stackLayout.set(currentPage);
                }
            } else {
                this.drawer.open();
            }

            // Handle websocket error messages
            ws.events.addListener("message-error", (msg) =>
                throwAlert({ message: msg, variant: "error" }),
            );

            // Updating the websocket server on store "server" canges
            store.addListener(
                "server",
                async (server) => {
                    ws.server = server;
                },
                true,
            );
        });
    }

    private initializeStore() {
        const store = this.store();
        const themeHandler = this.themeHandler();

        themeHandler.theme =
            store.getData("currentTheme")?.theme || themeHandler.theme;

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
        const store = this.store();
        const appBar = this.appBar();

        const stackLayout =
            this.shadowRoot!.querySelector<
                UIStackLayout<PicowStackLayoutPages>
            >(`ui-stack-layout`)!;

        stackLayout.register("devices", async () => {
            const page = new PicowDevicesPage();
            page.picowAppEvents = this.events;
            return page;
        });

        stackLayout.register("settings", async () => {
            return new PicowSettingsPage();
        });

        stackLayout.events.addListener("change", async ({ current }) => {
            const addItem = appBar.contentName("add")!;

            // Reset all layouts (AppBar buttons and title)
            appBar.title = "PicoW LED";
            addItem.hide();

            if (!current) {
                this.drawer.open();
                return;
            }

            store.setData("currentPage", current.name as PicowStackLayoutPages);
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
        const appBar = this.appBar();

        this.events.addListener("menu", () => this.drawer.open());

        const statusItem = appBar.contentName("status")!;
        const picowStatusLED =
            statusItem.querySelector<PicowStatusLED>(`picow-status-led`)!;

        const text = statusItem.querySelector(`ui-secondary`)!;

        ws.events.addListener("open", () => {
            picowStatusLED.active = true;
            text.innerHTML = "Online";
        });

        ws.events.addListener("close", () => {
            picowStatusLED.active = false;
            text.innerHTML = "Offline";
        });
    }

    private initializeDrawer() {
        // TODO: ...
    }
}
