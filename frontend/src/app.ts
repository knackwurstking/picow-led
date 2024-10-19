import "ui/lib/css/main.css";

import { css as CSS, html, LitElement, type PropertyValues } from "lit";
import { customElement } from "lit/decorators.js";
import { globalStylesToShadowRoot, UIAlerts, UIThemeHandler } from "ui";
import { PicowAppBar } from "./app/picow-app-bar";
import { PicowDrawer } from "./app/picow-drawer";
import { PicowStore } from "./types";

@customElement("picow-app")
export class PicowApp extends LitElement {
    private appBar: PicowAppBar = new PicowAppBar();
    private drawer: PicowDrawer = new PicowDrawer();
    private alerts: UIAlerts = new UIAlerts();

    private get store(): PicowStore {
        return document.querySelector<PicowStore>(`ui-store`)!;
    }

    private get themeHandler(): UIThemeHandler {
        return document.querySelector<UIThemeHandler>(`ui-theme-handler`)!;
    }

    static get styles() {
        return CSS`
            :host {
                position: fixed !important;
                top: 0;
                right: 0;
                bottom: 0;
                left: 0;
                overflow: auto;
            }
        `;
    }

    protected render() {
        return html`
            <ui-container style="width: 100%; height: 100%;">
                <ui-stack-layout></ui-stack-layout>
            </ui-container>

            ${this.appBar} ${this.drawer} ${this.alerts}
        `;
    }

    protected firstUpdated(_changedProperties: PropertyValues): void {
        globalStylesToShadowRoot(this.shadowRoot!);
    }
}
