import "ui/lib/css/main.css";

import { css as CSS, html, LitElement, type PropertyValues } from "lit";
import { customElement } from "lit/decorators.js";
import { globalStylesToShadowRoot, UIStackLayout, UIThemeHandler } from "ui";
import { PicowAppBar } from "./app/picow-app-bar";
import { PicowDrawer } from "./app/picow-drawer";
import { throwAlert } from "./lib/utils";
import { ws } from "./lib/websocket";
import { PicowStackLayoutPages, PicowStore } from "./types";

/**
 * **Tag**: picow-app
 */
@customElement("picow-app")
export class PicowApp extends LitElement {
    private appBar: PicowAppBar = new PicowAppBar();
    private drawer: PicowDrawer = new PicowDrawer();

    private store(): PicowStore {
        return document.querySelector<PicowStore>(`ui-store`)!;
    }

    private themeHandler(): UIThemeHandler {
        return document.querySelector<UIThemeHandler>(`ui-theme-handler`)!;
    }

    private stackLayout(): UIStackLayout<PicowStackLayoutPages> {
        return this.shadowRoot!.querySelector<
            UIStackLayout<PicowStackLayoutPages>
        >(`ui-stack-layout`)!;
    }

    static get styles() {
        return CSS`
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

            ${this.appBar} ${this.drawer}
        `;
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        globalStylesToShadowRoot(this.shadowRoot!);

        this.initializeStore();
        this.initializeStackLayout();
        this.initializeAppBar();

        setTimeout(() => {
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
                throwAlert({ message: msg, variant: "error" })
            );

            // Updating the websocket server on store "server" canges
            store.addListener(
                "server",
                async (server) => {
                    ws.server = server;
                },
                true
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
            true
        );
    }

    private initializeStackLayout() {
        const store = this.store();

        const stackLayout =
            this.shadowRoot!.querySelector<
                UIStackLayout<PicowStackLayoutPages>
            >(`ui-stack-layout`)!;

        stackLayout.register("devices", async () => {
            const page = new PicowDevicesPage();

            // TODO: Do app-bar stuff here

            return page;
        });

        stackLayout.register("settings", async () => {
            return new PicowSettingsPage();
        });

        stackLayout.events.addListener("change", ({ current }) => {
            // Reset all layouts (AppBar buttons and title)
            this.appBar.title = "PicoW LED";
            this.appBar.items.add.hide(); // TODO: Continue here...

            if (!current) {
                this.drawer.open();
                return;
            }

            store.setData("currentPage", current.name as PicowStackLayoutPages);
            switch (current.name) {
                case "devices":
                    store.setData("currentPage", current.name);
                    this.appBar.title = "Devices";
                    this.appBar.items.add.show();
                    break;

                case "settings":
                    store.setData("currentPage", current.name);
                    this.appBar.title = "Settings";
                    break;

                default:
                    this.appBar.title = current.name;
                    break;
            }
        });
    }

    private initializeAppBar() {
        this.appBar.events.addListener("menu", () => this.drawer.open());
    }
}
