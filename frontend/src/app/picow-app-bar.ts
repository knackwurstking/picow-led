import { html, LitElement, PropertyValues } from "lit";
import { customElement } from "lit/decorators.js";
import { Events, svg, UIAppBar, UIAppBarItem, UIIconButton } from "ui";
import { PicowStatusLED } from "../components/picow-status-led";
import { ws } from "../lib/websocket";
import { PicowAppBarEvents } from "../types";

/**
 * **Tag**: picow-app-bar
 *
 * **Public Methods**:
 *  - `title: string`
 *  - `root(): UIAppBar | null`
 */
@customElement("picow-app-bar")
export class PicowAppBar extends LitElement {
    public events: Events<PicowAppBarEvents> = new Events();

    protected render() {
        return html`
            <ui-app-bar position="top">
                <ui-app-bar-item name="menu" slot="left">
                    <ui-icon-button ghost>
                        ${svg.smoothieLineIcons.menu}
                    </ui-icon-button>
                </ui-app-bar-item>

                <ui-app-bar-item name="status" slot="left">
                    <ui-flex-grid-row align="flex-end" gap="0.25rem">
                        <picow-status-led></picow-status-led>
                        <ui-secondary
                            style="white-space: nowrap;"
                        ></ui-secondary>
                    </ui-flex-grid-row>
                </ui-app-bar-item>

                <ui-app-bar-item name="title" slot="center">
                    <h4 style="white-space: nowrap;">PicoW LED</h4>
                </ui-app-bar-item>

                <ui-app-bar-item name="add" slot="right">
                    <ui-icon-button ghost>
                        ${svg.smoothieLineIcons.plus}
                    </ui-icon-button>
                </ui-app-bar-item>
            </ui-app-bar>
        `;
    }

    protected updated(_changedProperties: PropertyValues): void {
        const root = this.root()!;

        let item: UIAppBarItem;

        // Forward "menu" button click event
        item = root.contentName("menu")!;
        {
            const button = item.contentAt<UIIconButton>(0);
            button.onclick = (ev) => this.events.dispatch("menu", ev);
        }

        // Initialize the satus led
        const statusItem = root.contentName("status")!;
        {
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

        // Forward "add" button click
        item = root.contentName("add")!;
        {
            const button = item.contentAt<UIIconButton>(0);
            button.onclick = (ev) => this.events.dispatch("add", ev);
        }
    }

    public get title() {
        const root = this.root();
        if (!root) return "";
        return root.contentName("title")!.contentAt(0).innerText;
    }

    public set title(value: string) {
        const root = this.root()!;
        root.contentName("title")!.contentAt(0).innerText = value;
    }

    public root(): UIAppBar | null {
        return this.shadowRoot?.querySelector<UIAppBar>(`ui-app-bar`)! || null;
    }
}
