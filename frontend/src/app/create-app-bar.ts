import { plus as svgAdd, menu as svgMenu } from "ui/svg/smoothie-line-icons";

import {
    Events,
    UIAppBar,
    UIAppBarItem,
    UIIconButton,
    UISecondary,
    html,
} from "ui";
import type StatusLED from "../components/status-led";
import ws from "../lib/websocket";
import { styles } from "ui";

export interface AppBar {
    element: import("ui").UIAppBar;
    events: import("ui").Events<AppBar_Events>;
    items: {
        menu: import("ui").UIAppBarItem<import("ui").UIIconButton>;
        title: import("ui").UIAppBarItem<HTMLElement>;
        add: import("ui").UIAppBarItem<import("ui").UIIconButton>;
    };
    buttons: {
        menu: import("ui").UIIconButton;
        add: import("ui").UIIconButton;
    };
    title: string;
}

export interface AppBar_Events {
    menu: MouseEvent & { currentTarget: import("ui").UIIconButton };
    add: MouseEvent & { currentTarget: import("ui").UIIconButton };
}

export default async function (): Promise<AppBar> {
    const el = new UIAppBar();
    el.ui.position = "top";

    el.innerHTML = html`
        <ui-app-bar-item name="menu" slot="left">
            <ui-icon-button ghost> ${svgMenu} </ui-icon-button>
        </ui-app-bar-item>

        <ui-app-bar-item name="status" slot="left">
            <ui-flex-grid-row align="flex-end" gap="0.25rem">
                <status-led></status-led>
                <ui-secondary style="white-space: nowrap;"></ui-secondary>
            </ui-flex-grid-row>
        </ui-app-bar-item>

        <ui-app-bar-item name="title" slot="center">
            <h4 style="white-space: nowrap;">PicoW LED</h4>
        </ui-app-bar-item>

        <ui-app-bar-item name="add" slot="right">
            <ui-icon-button ghost> ${svgAdd} </ui-icon-button>
        </ui-app-bar-item>
    `;

    const events = new Events<AppBar_Events>();

    //
    // Left Slot: Menu
    //

    const menu = el.querySelector<UIAppBarItem<UIIconButton>>(
        `ui-app-bar-item[name="menu"]`
    );

    menu.ui.child.ui.events.on("click", (ev) => {
        events.dispatch("menu", ev);
    });

    //
    // Left Slot: Status
    //

    {
        const status = el.querySelector<UIAppBarItem<HTMLDivElement>>(
            `ui-app-bar-item[name="status"]`
        );

        const led = status.ui.child.querySelector<StatusLED>("status-led");
        const statusText =
            status.ui.child.querySelector<UISecondary>(`ui-secondary`);

        ws.events.on("open", () => {
            led.setAttribute("active", "");
            statusText.innerText = "Online";
        });

        ws.events.on("close", () => {
            led.removeAttribute("active");
            statusText.innerText = "Offline";
        });
    }

    //
    // Center Slot: Title
    //

    const title = el.querySelector<UIAppBarItem<HTMLElement>>(
        `ui-app-bar-item[name="title"]`
    );

    //
    // Right Slot: Add
    //

    const add = el.querySelector<UIAppBarItem<UIIconButton>>(
        `ui-app-bar-item[name="add"]`
    );

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
