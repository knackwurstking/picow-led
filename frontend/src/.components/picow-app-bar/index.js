import svgMenu from "ui/src/svg/smoothie-line-icons/menu";
import svgAdd from "ui/src/svg/smoothie-line-icons/plus";

import { Events, html, UIAppBar } from "ui";

export class PicowAppBar extends UIAppBar {
    static register = () => {
        customElements.define("picow-app-bar", PicowAppBar);
    };

    constructor() {
        super();

        this.picow = {
            root: this,

            /** @type {Events<PicowAppBar_Events>} */
            events: new Events(),

            get title() {
                return this.root.items.title.ui.child.innerText;
            },

            set title(v) {
                this.root.items.title.ui.child.innerText = v || "";
            },
        };

        this.items;

        this.render();
    }

    render() {
        this.innerHTML = html`
            <ui-app-bar-item name="menu" slot="left">
                <ui-icon-button ghost> ${svgMenu} </ui-icon-button>
            </ui-app-bar-item>

            <ui-app-bar-item name="title" slot="center">
                <h4 style="white-space: nowrap;"></h4>
            </ui-app-bar-item>

            <ui-app-bar-item name="add" slot="right">
                <ui-icon-button ghost> ${svgAdd} </ui-icon-button>
            </ui-app-bar-item>
        `;

        this.items = {
            menu: this.createMenuItem(),

            title: this.createTitleItem(),

            add: this.createAddItem(),
        };
    }

    /** @private */
    createMenuItem() {
        /** @type {import("ui").UIAppBarItem<import("ui").UIIconButton>} */
        const item = this.querySelector(`ui-app-bar-item[name="menu"]`);

        item.ui.child.ui.events.on("click", (data) => {
            this.picow.events.dispatch("menu", data);
        });

        return item;
    }

    /** @private */
    createTitleItem() {
        /** @type {import("ui").UIAppBarItem<HTMLElement>} */
        const item = this.querySelector(`ui-app-bar-item[name="title"]`);
        return item;
    }

    /** @private */
    createAddItem() {
        /** @type {import("ui").UIAppBarItem<import("ui").UIIconButton>} */
        const item = this.querySelector(`ui-app-bar-item[name="add"]`);

        item.ui.child.ui.events.on("click", (data) => {
            this.picow.events.dispatch("add", data);
        });

        return item;
    }
}
