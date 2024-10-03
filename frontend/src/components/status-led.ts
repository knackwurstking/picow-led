import { html } from "ui";

/**
 * HTML: `status-led`
 *
 * Attributes:
 *   - __active__: *boolean*
 */
export default class StatusLED extends HTMLElement {
    constructor() {
        super();
        this.#renderStatusLED();
    }

    #renderStatusLED() {
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

console.debug(`Register the "status-led"`);
customElements.define("status-led", StatusLED);
