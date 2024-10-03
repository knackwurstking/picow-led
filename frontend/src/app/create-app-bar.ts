import { plus as svgAdd, menu as svgMenu } from "ui/svg/smoothie-line-icons";

import { Events, UIAppBar, UIAppBarItem, UIIconButton, html } from "ui";

interface AppBar {
    element: UIAppBar;
    events: Events<PicowAppBar_Events>;
    items: {
        menu: UIAppBarItem<UIIconButton>;
        title: UIAppBarItem<HTMLElement>;
        add: UIAppBarItem<UIIconButton>;
    };
    buttons: {
        menu: UIIconButton;
        add: UIIconButton;
    };
    title: string;
}

export default async function (): Promise<AppBar> {
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

    const events = new Events<PicowAppBar_Events>();

    const menu = el.querySelector(
        `ui-app-bar-item[name="menu"]`
    ) as UIAppBarItem<UIIconButton>;

    menu.ui.child.ui.events.on("click", (ev) => {
        events.dispatch("menu", ev);
    });

    const title = el.querySelector(
        `ui-app-bar-item[name="title"]`
    ) as UIAppBarItem<HTMLElement>;

    const add = el.querySelector(
        `ui-app-bar-item[name="add"]`
    ) as UIAppBarItem<UIIconButton>;

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
