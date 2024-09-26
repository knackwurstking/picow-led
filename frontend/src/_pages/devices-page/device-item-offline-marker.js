import { html } from "ui";

export class DeviceItemOfflineMarker extends HTMLElement {
    static register = () => {
        customElements.define(
            "device-item-offline-marker",
            DeviceItemOfflineMarker
        );
    };

    constructor() {
        super();

        this.picow = {
            root: this,

            get hide() {
                return this.root.hasAttribute("hide");
            },

            set hide(state) {
                if (!state) {
                    this.root.removeAttribute("hide");
                    return;
                }

                this.root.setAttribute("hide", "");
            },
        };

        this.render();
    }

    render() {
        this.attachShadow({ mode: "open" });
        this.shadowRoot.innerHTML = html`
            <style>
                :host {
                    padding-left: var(--ui-spacing);
                    display: block;
                    position: absolute !important;
                    top: -0.25rem;
                    left: 50%;
                    color: var(--ui-destructive);
                    transform: translateX(-50%);
                }

                :host([hide]) {
                    display: none;
                }
            </style>

            <slot></slot>
        `;

        this.innerHTML = html`
            <ui-secondary style="text-wrap: nowrap;"
                >Device Marked Offline</ui-secondary
            >
        `;
    }
}
