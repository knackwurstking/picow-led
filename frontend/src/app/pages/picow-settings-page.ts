import { css, CSSResult, html, TemplateResult } from "lit";
import { customElement } from "lit/decorators.js";
import {
    UICheck,
    UIInput,
    UISelect,
    UIStackLayoutPage,
    UIThemeHandler,
    UIThemeHandlerTheme,
} from "ui";
import { PicowStore } from "../../types";

@customElement("picow-settings-page")
export class PicowSettingsPage extends UIStackLayoutPage {
    name = "settings";

    private store: PicowStore = document.querySelector(`ui-store`)!;
    private themeHandler: UIThemeHandler =
        document.querySelector(`ui-theme-handler`)!;

    static get styles(): CSSResult {
        return css`
            ${UIStackLayoutPage.styles}

            :host {
                padding-top: var(--ui-app-bar-height);
                overflow: auto;
            }
        `;
    }

    protected render(): TemplateResult<1> {
        let timeout: NodeJS.Timeout | null = null;
        const timeoutValueMS: number = 250;

        const resetTimeoutHandler = () => {
            if (!timeout) return;
            clearTimeout(timeout);
            timeout = null;
        };

        return html`
            <ui-flex-grid gap="0.25rem">
                <ui-flex-grid-item>
                    <ui-label primary="Use SSL connections" ripple>
                        <ui-check
                            name="ssl"
                            ?checked=${this.store.getData("server")?.ssl}
                            @input=${async (ev: Event) => {
                                resetTimeoutHandler();
                                const target = ev.currentTarget as UICheck;

                                timeout = setTimeout(() => {
                                    this.store.updateData(
                                        "server",
                                        (server) => {
                                            server.ssl = target.checked;
                                            return server;
                                        },
                                    );
                                }, timeoutValueMS);
                            }}
                        ></ui-check>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Host">
                        <ui-input
                            name="host"
                            value="${this.store.getData("server")?.host}"
                            @input=${async (ev: Event) => {
                                resetTimeoutHandler();
                                const target = ev.currentTarget as UIInput;

                                timeout = setTimeout(() => {
                                    this.store.updateData(
                                        "server",
                                        (server) => {
                                            server.host = target.value;
                                            return server;
                                        },
                                    );
                                }, timeoutValueMS);
                            }}
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Server Port">
                        <ui-input
                            name="port"
                            type="number"
                            value="${this.store.getData("server")?.port}"
                            @input=${async (ev: Event) => {
                                resetTimeoutHandler();
                                const target = ev.currentTarget as UIInput;

                                timeout = setTimeout(() => {
                                    this.store.updateData(
                                        "server",
                                        (server) => {
                                            server.port = target.value;
                                            return server;
                                        },
                                    );
                                }, timeoutValueMS);
                            }}
                        ></ui-input>
                    </ui-label>
                </ui-flex-grid-item>

                <ui-flex-grid-item>
                    <ui-label primary="Theme">
                        <ui-select
                            name="theme"
                            keep-open
                            @change=${(ev: Event) => {
                                const target = ev.currentTarget as UISelect;

                                const option = target.selected();
                                if (option === null) return;

                                this.themeHandler.theme = target.selected()
                                    ?.value as UIThemeHandlerTheme;

                                this.store.setData("currentTheme", {
                                    theme: this.themeHandler.theme,
                                });
                            }}
                        >
                            <ui-select-option
                                value="original"
                                ?selected=${this.themeHandler.theme ===
                                "original"}
                            >
                                Original
                            </ui-select-option>

                            <ui-select-option
                                value="gruvbox"
                                ?selected=${this.themeHandler.theme ===
                                "gruvbox"}
                            >
                                Gruvbox
                            </ui-select-option>
                        </ui-select>
                    </ui-label>
                </ui-flex-grid-item>
            </ui-flex-grid>
        `;
    }
}
