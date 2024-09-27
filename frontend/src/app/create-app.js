import { html, styles } from "ui";
import createAppBar from "./create-app-bar";
import createDrawer from "./create-drawer";

/**
 * @typedef {{
 *  element: HTMLDivElement;
 * }} App
 */

/**
 * @returns {App}
 */
export default function () {
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

    const appBar = createAppBar();
    el.querySelector(`div.app-bar`).replaceWith(appBar.element);

    const drawer = createDrawer();
    el.querySelector(`div.drawer`).replaceWith(drawer.element);

    appBar.buttons.menu.ui.events.on("click", () => drawer.open());

    // TODO: Initialize the store and set "server" to `deviceEvents` and `devicesEvents`
    // TODO: Initialize the layout
    //          - Register pages
    //          - Handle the stack layout "change" event and setup drawer items, title, ...

    return {
        element: el,
    };
}
