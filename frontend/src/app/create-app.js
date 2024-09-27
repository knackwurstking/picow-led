import { html, styles } from "ui";
import createAppBar from "./create-app-bar";
import createDrawer from "./create-drawer";

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

    // TODO: Initialize the layout
    //          - Register pages
    //          - Handle the stack layout "change" event and setup drawer items, title, ...

    // ----------------- //
    // Create the AppBar //
    // ----------------- //

    const appBar = await createAppBar();
    el.querySelector(`div.app-bar`).replaceWith(appBar.element);

    // ----------------- //
    // Create the Drawer //
    // ----------------- //

    const drawer = await createDrawer();
    el.querySelector(`div.drawer`).replaceWith(drawer.element);

    appBar.buttons.menu.ui.events.on("click", () => drawer.open());

    return {
        element: el,
    };
}
