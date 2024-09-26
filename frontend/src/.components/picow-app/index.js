import { CleanUp, html, styles } from "ui";
import { deviceEvents, devicesEvents } from "../../lib";
import { DevicesPage, SettingsPage } from "../../pages";
import { PicowAppBar } from "../picow-app-bar";
import { PicowDrawer } from "../picow-drawer";

export class PicowApp extends HTMLElement {
    static register = () => {
        PicowAppBar.register();
        PicowDrawer.register();

        DevicesPage.register();
        SettingsPage.register();

        customElements.define("picow-app", PicowApp);
    };

    constructor() {
        super();

        this.cleanup = new CleanUp();

        /** @type {PicowStore} */
        this.uiStore;
        /** @type {PicowStackLayout} */
        this.uiLayout;
        /** @type {import("../picow-app-bar").PicowAppBar} */
        this.picowAppBar;
        /** @type {import("../picow-drawer").PicowDrawer} */
        this.picowDrawer;

        this.render();
    }

    render() {
        this.attachShadow({ mode: "open" });
        this.shadowRoot.innerHTML = html`
            <style>
                :host {
                    display: block;
                }
            </style>

            <slot></slot>
        `;

        this.innerHTML = html`
            <ui-theme-handler mode="dark"></ui-theme-handler>
            <ui-store storageprefix="picow:" storage></ui-store>

            <picow-app-bar position="top"></picow-app-bar>
            <picow-drawer></picow-drawer>
            <ui-container
                style="${styles({
                    width: "100%",
                    height: "100%",
                })}"
            >
                <ui-stack-layout></ui-stack-layout>
            </ui-container>

            <ui-alerts></ui-alerts>
        `;

        this.uiStore = this.querySelector("ui-store");
        this.uiLayout = this.querySelector("ui-stack-layout");
        this.picowAppBar = this.querySelector("picow-app-bar");
        this.picowDrawer = this.querySelector("picow-drawer");

        this.createStore();
        this.createLayout();
        this.createAppBar();
    }

    connectedCallback() {
        this.resetLayout();

        if (!!this.uiStore.ui.get("currentPage")) {
            this.uiLayout.ui.set(this.uiStore.ui.get("currentPage"));
        } else {
            this.picowDrawer.ui.open = true;
        }

        this.cleanup.add(
            this.uiStore.ui.on(
                "server",
                async (server) => {
                    deviceEvents.server = server;
                    devicesEvents.server = server;
                },
                true
            )
        );
    }

    disconnectedCallback() {
        this.cleanup.run();
        deviceEvents.close();
        devicesEvents.close();
    }

    /** @private */
    createStore() {
        this.uiStore.ui.set("devices", [], true);
        this.uiStore.ui.set("currentPage", null, true);
        this.uiStore.ui.set(
            "server",
            {
                ssl: !!location.protocol.match(/(https)/),
                host: location.hostname,
                port: location.port,
            },
            true
        );

        if (!Object.hasOwn(this.uiStore.ui.get("server"), "ssl")) {
            this.uiStore.ui.update("server", (server) => {
                // @ts-ignore
                if (Object.hasOwn(server, "protocol")) delete server.protocol;

                return {
                    ...server,
                    ssl: !!location.protocol.match(/(https)/),
                };
            });
        }
    }

    /** @private */
    createLayout() {
        this.uiLayout.ui.register("devices", () => {
            return new DevicesPage();
        });

        this.uiLayout.ui.register("settings", () => {
            return new SettingsPage();
        });

        this.uiLayout.ui.events.on("change", ({ newPage }) => {
            this.resetLayout();

            if (!newPage) {
                this.picowDrawer.ui.open = true;
                return;
            }

            switch (newPage.ui.name) {
                case "devices":
                    this.uiStore.ui.set("currentPage", newPage.ui.name);
                    this.picowAppBar.picow.title = "Devices";
                    this.picowAppBar.items.add.ui.show();
                    break;

                case "settings":
                    this.uiStore.ui.set("currentPage", newPage.ui.name);
                    this.picowAppBar.picow.title = "Settings";
                    break;
            }
        });
    }

    /** @private */
    createAppBar() {
        this.picowAppBar.picow.events.on("menu", () => {
            this.picowDrawer.ui.open = !this.picowDrawer.ui.open;
        });
    }

    /** @private */
    resetLayout() {
        this.picowAppBar.picow.title = "PicoW LED";
        this.picowAppBar.items.add.ui.hide();
    }
}
