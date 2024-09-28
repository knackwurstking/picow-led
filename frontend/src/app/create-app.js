import { html, styles } from "ui";
import { deviceEvents, devicesEvents } from "../lib";
import createAppBar from "./create-app-bar";
import createDrawer from "./create-drawer";
import PicowDevicesPage from "./pages/picow-devices-page";
import PicowSettingsPage from "./pages/picow-settings-page";

/**
 * @typedef {{
 *  element: HTMLDivElement;
 * }} App
 */

/**
 * @returns {Promise<App>}
 */
export default async function () {
    const el = document.createElement("div");

    el.style.width = "100%";
    el.style.height = "100%";

    el.innerHTML = html`
        <ui-theme-handler mode="dark"></ui-theme-handler>
        <ui-store storageprefix="picow:" storage></ui-store>

        <div class="app-bar"></div>
        <div class="drawer"></div>

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

    // ---------------- //
    // Create the Store //
    // ---------------- //

    /**
     * @type {PicowStore}
     */
    const store = el.querySelector(`ui-store`);
    store.ui.set("devices", [], true);
    store.ui.set("currentPage", null, true);
    store.ui.set(
        "server",
        {
            ssl: !!location.protocol.match(/(https)/),
            host: location.hostname,
            port: location.port,
        },
        true
    );

    if (!Object.hasOwn(store.ui.get("server"), "ssl")) {
        store.ui.update("server", (server) => {
            // @ts-ignore
            if (Object.hasOwn(server, "protocol")) delete server.protocol;

            return {
                ...server,
                ssl: !!location.protocol.match(/(https)/),
            };
        });
    }

    // ---------------------- //
    // Create the StackLayout //
    // ---------------------- //

    /**
     * @type {import("ui").UIStackLayout<PicowStackLayout_Pages>}
     */
    const stackLayout = el.querySelector(`ui-stack-layout`);

    stackLayout.ui.register("devices", async () => {
        return new PicowDevicesPage({ store, appBar });
    });

    stackLayout.ui.register("settings", async () => {
        return new PicowSettingsPage({ store });
    });

    stackLayout.ui.events.on("change", ({ newPage }) => {
        // Reset all layouts (AppBar buttons and title)
        appBar.title = "PicoW LED";
        appBar.items.add.ui.hide();

        if (!newPage) {
            drawer.open();
            return;
        }

        switch (newPage.ui.name) {
            case "devices":
                store.ui.set("currentPage", newPage.ui.name);
                appBar.title = "Devices";
                appBar.items.add.ui.show();
                break;

            case "settings":
                store.ui.set("currentPage", newPage.ui.name);
                appBar.title = "Settings";
                break;

            default:
                appBar.title = newPage.ui.name;
                break;
        }
    });

    // ----------------- //
    // Create the AppBar //
    // ----------------- //

    const appBar = await createAppBar();
    el.querySelector(`div.app-bar`).replaceWith(appBar.element);

    appBar.events.on("menu", () => drawer.open());

    // ----------------- //
    // Create the Drawer //
    // ----------------- //

    const drawer = await createDrawer({ stackLayout });
    el.querySelector(`div.drawer`).replaceWith(drawer.element);

    // ---------------- //
    // Lets get started //
    // ---------------- //

    setTimeout(() => {
        // Set the start page
        if (!!store.ui.get("currentPage")) {
            stackLayout.ui.set(store.ui.get("currentPage"));
        } else {
            drawer.open();
        }

        // Handler server changes
        store.ui.on(
            "server",
            async (server) => {
                deviceEvents.server = server;
                devicesEvents.server = server;
            },
            true
        );
    });

    return {
        element: el,
    };
}
