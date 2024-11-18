import { customElement, property } from "lit/decorators.js";

import { css, html, LitElement } from "lit";

@customElement("picow-status-led")
class PicowStatusLED extends LitElement {
    @property({ type: Boolean, attribute: "active", reflect: true })
    active: boolean = false;

    static get styles() {
        return css`
            :host {
                display: block;
            }

            :host .outer {
                padding: 0.25rem;
                border-radius: 50%;
                background-color: transparent;
            }

            :host .inner {
                width: 1rem;
                height: 1rem;
                border-radius: 50%;
                filter: blur(2px);

                transition: background-color 0.15s linear;
            }

            :host(:not([active])) .inner {
                background-color: rgb(255, 0, 0);
            }

            :host([active]) .inner {
                background-color: rgb(0, 255, 0);
            }
        `;
    }

    protected render() {
        return html`
            <div class="outer">
                <div class="inner"></div>
            </div>
        `;
    }
}

export default PicowStatusLED;
