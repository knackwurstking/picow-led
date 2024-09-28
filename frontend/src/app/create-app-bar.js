import { menu as svgMenu, plus as svgAdd } from "ui/svg/smoothie-line-icons";
import { Events, UIAppBar, html } from "ui";

/**
 * @typedef {{
 *  element: UIAppBar;
 *  events: Events<PicowAppBar_Events>;
 *  items: {
 *      menu: import("ui").UIAppBarItem<import("ui").UIIconButton>;
 *      title: import("ui").UIAppBarItem<HTMLElement>;
 *      add: import("ui").UIAppBarItem<import("ui").UIIconButton>;
 *  };
 *  buttons: {
 *      menu: import("ui").UIIconButton;
 *      add: import("ui").UIIconButton;
 *  };
 *  title: string;
 * }} AppBar
 */

/**
 * @returns {Promise<AppBar>}
 */
export default async function () {
    const el = new UIAppBar();
    el.ui.position = "top";

    el.innerHTML = html`
        <ui-app-bar-item name="menu" slot="left">
            <ui-icon-button ghost> ${svgMenu} </ui-icon-button>
        </ui-app-bar-item>

        <ui-app-bar-item name="title" slot="center">
            <h4 style="white-space: nowrap;">PicoW LED</h4>
        </ui-app-bar-item>

        <ui-app-bar-item name="add" slot="right">
            <ui-icon-button ghost> ${svgAdd} </ui-icon-button>
        </ui-app-bar-item>
    `;

    /** @type {Events<PicowAppBar_Events>} */
    const events = new Events();

    /**
     * @type {import("ui").UIAppBarItem<import("ui").UIIconButton>}
     */
    const menu = el.querySelector(`ui-app-bar-item[name="menu"]`);
    menu.ui.child.ui.events.on("click", (ev) => {
        events.dispatch("menu", ev);
    });

    /**
     * @type {import("ui").UIAppBarItem<HTMLElement>}
     */
    const title = el.querySelector(`ui-app-bar-item[name="title"]`);

    /**
     * @type {import("ui").UIAppBarItem<import("ui").UIIconButton>}
     */
    const add = el.querySelector(`ui-app-bar-item[name="add"]`);
    add.ui.child.ui.events.on("click", (ev) => {
        events.dispatch("add", ev);
    });
    add.ui.hide();

    return {
        element: el,
        events,

        get items() {
            return { menu, title, add };
        },

        get buttons() {
            return { menu: menu.ui.child, add: add.ui.child };
        },

        get title() {
            return title.ui.child.innerText;
        },

        set title(value) {
            title.ui.child.innerText = value || "";
        },
    };
}
