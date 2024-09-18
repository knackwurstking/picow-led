import { html } from "ui";

/**
 * HTML: `status-led`
 *
 * Attributes:
 *   - active: boolean
 */
export class StatusLED extends HTMLElement {
    static register = () => {
        if (!customElements.get("status-led")) {
            customElements.define("status-led", StatusLED);
        }
    };

    constructor() {
        super();
        this.shadowRender();
    }

    shadowRender() {
        this.attachShadow({ mode: "open" });
        this.shadowRoot.innerHTML = html`
            <style>
                :host {
                    display: block;
                }

                :host .outer {
                    padding: 0.25rem;
                    border-radius: 50%;
                    background-color: var(--ui-bg);
                }

                :host .inner {
                    width: 1rem;
                    height: 1rem;
                    border-radius: 50%;
                    filter: blur(2px);
                }

                :host(:not([active])) .inner {
                    background-color: rgb(255, 0, 0);
                }

                :host([active]) .inner {
                    background-color: rgb(0, 255, 0);
                }
            </style>

            <div class="outer">
                <div class="inner"></div>
            </div>
        `;
    }
}
