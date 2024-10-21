import { html, LitElement, PropertyValues } from "lit";
import { customElement } from "lit/decorators.js";
import { UIButton, UIDrawer, UIStackLayout } from "ui";
import { PicowStackLayoutPages } from "../types";

/**
 * **App**: picow-drawer
 *
 * **Public Methods**:
 *  - `root(): UIDrawer`
 *  - `open()`
 *  - `close()`
 */
@customElement("picow-drawer")
export class PicowDrawer extends LitElement {
    private stackLayout(): UIStackLayout<PicowStackLayoutPages> {
        return this.shadowRoot!.querySelector<
            UIStackLayout<PicowStackLayoutPages>
        >(`ui-stack-layout`)!;
    }

    protected render() {
        return html`
            <ui-drawer>
                <ui-drawer-group no-fold>
                    <ui-drawer-group-item>
                        <ui-button
                            name="devices"
                            style="justify-content: flex-start;"
                            variant="ghost"
                            color="primary"
                        >
                            Devices
                        </ui-button>
                    </ui-drawer-group-item>

                    <ui-drawer-group-item>
                        <ui-button
                            name="groups"
                            style="justify-content: flex-start;"
                            variant="ghost"
                            color="primary"
                            disabled
                        >
                            Groups
                        </ui-button>
                    </ui-drawer-group-item>

                    <ui-drawer-group-item>
                        <ui-button
                            name="scenes"
                            style="justify-content: flex-start;"
                            variant="ghost"
                            color="primary"
                            disabled
                        >
                            Scenes
                        </ui-button>
                    </ui-drawer-group-item>

                    <ui-drawer-group-item>
                        <ui-button
                            name="settings"
                            style="justify-content: flex-start;"
                            variant="ghost"
                            color="primary"
                        >
                            Settings
                        </ui-button>
                    </ui-drawer-group-item>
                </ui-drawer-group>
            </ui-drawer>
        `;
    }

    protected updated(_changedProperties: PropertyValues): void {
        const root = this.root();
        const stackLayout = this.stackLayout();
        let button: UIButton;

        const pages: PicowStackLayoutPages[] = ["devices", "settings"];
        for (const name of pages) {
            button = this.shadowRoot!.querySelector(
                `ui-button[name="devices"]`,
            )!;
            button.onclick = () => {
                stackLayout.clear();
                stackLayout.set(name);
                root.open = false;
            };
        }
    }

    public root(): UIDrawer {
        return this.shadowRoot!.querySelector(`ui-drawer`)!;
    }

    public open() {
        this.root().open = true;
    }

    public close() {
        this.root().open = false;
    }
}
